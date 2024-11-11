package plugin_mp4

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"m7s.live/pro/plugin/mp4/pb"

	m7s "m7s.live/pro"
	"m7s.live/pro/pkg/util"
	mp4 "m7s.live/pro/plugin/mp4/pkg"
	"m7s.live/pro/plugin/mp4/pkg/box"
)

type ContentPart struct {
	*os.File
	Start int64
	Size  int
}

func (p *MP4Plugin) List(ctx context.Context, req *pb.ReqRecordList) (resp *pb.ResponseList, err error) {
	var streams []m7s.RecordStream
	if p.DB == nil {
		err = fmt.Errorf("db not init")
		return
	}
	startTime, endTime, err := util.TimeRangeQueryParse(url.Values{"range": []string{req.Range}, "start": []string{req.Start}, "end": []string{req.End}})
	if err != nil {
		return
	}
	if req.StreamPath == "" {
		p.DB.Find(&streams, "end_time>? AND start_time<?", startTime, endTime)
	} else if strings.Contains(req.StreamPath, "*") {
		p.DB.Find(&streams, "end_time>? AND start_time<? AND stream_path like ?", startTime, endTime, strings.ReplaceAll(req.StreamPath, "*", "%"))
	} else {
		p.DB.Find(&streams, "end_time>? AND start_time<? AND stream_path=?", startTime, endTime, req.StreamPath)
	}
	resp = &pb.ResponseList{}
	for _, stream := range streams {
		resp.Data = append(resp.Data, &pb.RecordFile{
			Id:         uint32(stream.ID),
			StartTime:  timestamppb.New(stream.StartTime),
			EndTime:    timestamppb.New(stream.EndTime),
			FilePath:   stream.FilePath,
			StreamPath: stream.StreamPath,
		})
	}
	return
}

func (p *MP4Plugin) Catalog(ctx context.Context, req *emptypb.Empty) (resp *pb.ResponseCatalog, err error) {
	resp = &pb.ResponseCatalog{}
	var result []struct {
		StreamPath string
		Count      uint
		StartTime  time.Time
		EndTime    time.Time
	}
	err = p.DB.Model(&m7s.RecordStream{}).Select("stream_path,count(id) as count,min(start_time) as start_time,max(end_time) as end_time").Group("stream_path").Find(&result).Error
	if err != nil {
		return
	}
	for _, row := range result {
		resp.Data = append(resp.Data, &pb.Catalog{
			StreamPath: row.StreamPath,
			Count:      uint32(row.Count),
			StartTime:  timestamppb.New(row.StartTime),
			EndTime:    timestamppb.New(row.EndTime),
		})
	}
	return
}

func (p *MP4Plugin) Delete(ctx context.Context, req *pb.ReqRecordDelete) (resp *pb.ResponseDelete, err error) {
	return
}

func (p *MP4Plugin) download(w http.ResponseWriter, r *http.Request) {
	streamPath := r.PathValue("streamPath")
	startTime, endTime, err := util.TimeRangeQueryParse(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.Info("download", "streamPath", streamPath, "start", startTime, "end", endTime)
	var streams []m7s.RecordStream
	p.DB.Find(&streams, "end_time>? AND start_time<? AND stream_path=?", startTime, endTime, streamPath)
	muxer := mp4.NewMuxer(0)
	var n int
	n, err = w.Write(box.MakeFtypBox(box.TypeISOM, 0x200, box.TypeISOM, box.TypeISO2, box.TypeAVC1, box.TypeMP41))
	if err != nil {
		return
	}
	muxer.CurrentOffset = int64(n)
	var lastTs, tsOffset int64
	var parts []*ContentPart
	sampleOffset := muxer.CurrentOffset + box.BasicBoxLen*2
	mdatOffset := sampleOffset
	var audioTrack, videoTrack *mp4.Track
	var file *os.File
	streamCount := len(streams)
	for i, stream := range streams {
		tsOffset = lastTs
		file, err = os.Open(stream.FilePath)
		if err != nil {
			return
		}
		p.Info("read", "file", file.Name())
		demuxer := mp4.NewDemuxer(file)
		err = demuxer.Demux()
		if err != nil {
			return
		}
		if i == 0 {
			for _, track := range demuxer.Tracks {
				t := muxer.AddTrack(track.Cid)
				t.ExtraData = track.ExtraData
				if track.Cid.IsAudio() {
					audioTrack = t
					t.SampleSize = track.SampleSize
					t.SampleRate = track.SampleRate
					t.ChannelCount = track.ChannelCount
				} else if track.Cid.IsVideo() {
					videoTrack = t
					t.Width = track.Width
					t.Height = track.Height
				}
			}
			startTimestamp := startTime.Sub(stream.StartTime).Milliseconds()
			var startSample *box.Sample
			if startSample, err = demuxer.SeekTime(uint64(startTimestamp)); err != nil {
				tsOffset = 0
				continue
			}
			tsOffset = -int64(startSample.DTS)
		}
		var part *ContentPart
		for track, sample := range demuxer.RangeSample {
			if i == streamCount-1 && int64(sample.DTS) > endTime.Sub(stream.StartTime).Milliseconds() {
				break
			}
			if part == nil {
				part = &ContentPart{
					File:  file,
					Start: sample.Offset,
				}
			}
			part.Size += sample.Size
			lastTs = int64(sample.DTS + uint64(tsOffset))
			fixSample := *sample
			fixSample.DTS += uint64(tsOffset)
			fixSample.PTS += uint64(tsOffset)
			fixSample.Offset += sampleOffset - part.Start
			if track.Cid.IsAudio() {
				audioTrack.AddSampleEntry(fixSample)
			} else if track.Cid.IsVideo() {
				videoTrack.AddSampleEntry(fixSample)
			}
		}
		if part != nil {
			sampleOffset += int64(part.Size)
			parts = append(parts, part)
		}
	}
	moovSize := muxer.GetMoovSize()
	for _, track := range muxer.Tracks {
		for i := range track.Samplelist {
			track.Samplelist[i].Offset += int64(moovSize)
		}
	}
	err = muxer.WriteMoov(w)
	if err != nil {
		return
	}
	var mdatBox = box.MediaDataBox(sampleOffset - mdatOffset)
	boxLen, buf := mdatBox.Encode()
	if boxLen == box.BasicBoxLen*2 {
		w.Write(buf)
	} else {
		freeBox := box.NewBasicBox(box.TypeFREE)
		freeBox.Size = box.BasicBoxLen
		_, free := freeBox.Encode()
		w.Write(free)
		w.Write(buf)
	}
	var written, totalWritten int64
	for _, part := range parts {
		part.Seek(part.Start, io.SeekStart)
		written, err = io.CopyN(w, part.File, int64(part.Size))
		if err != nil {
			return
		}
		totalWritten += written
		part.Close()
	}
}
