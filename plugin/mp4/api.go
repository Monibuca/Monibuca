package plugin_mp4

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"m7s.live/v5/pkg/config"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	m7s "m7s.live/v5"
	"m7s.live/v5/pkg"
	"m7s.live/v5/pkg/util"
	"m7s.live/v5/plugin/mp4/pb"
	mp4 "m7s.live/v5/plugin/mp4/pkg"
	"m7s.live/v5/plugin/mp4/pkg/box"
)

type ContentPart struct {
	*os.File
	Start int64
	Size  int
}

func (p *MP4Plugin) List(ctx context.Context, req *pb.ReqRecordList) (resp *pb.ResponseList, err error) {
	var streams []m7s.RecordStream
	if p.DB == nil {
		err = pkg.ErrNoDB
		return
	}

	offset := (req.PageNum - 1) * req.PageSize // 计算偏移量
	var totalCount int64                       //总条数

	// 查询总记录数
	countQuery := p.DB.Model(&m7s.RecordStream{})
	// 查询当前页的数据
	query := p.DB.Model(&m7s.RecordStream{})
	if req.PageSize > 0 {
		query = query.Limit(int(req.PageSize)).Offset(int(offset))
	}
	startTime, endTime, err := util.TimeRangeQueryParse(url.Values{"range": []string{req.Range}, "start": []string{req.Start}, "end": []string{req.End}})
	if err != nil {
		return
	}
	var condition string = "end_time>? AND start_time<?"
	var values []any = []any{startTime, endTime}
	if strings.Contains(req.StreamPath, "*") {
		condition += " AND stream_path like ?"
		values = append(values, strings.ReplaceAll(req.StreamPath, "*", "%"))
	} else if req.StreamPath != "" {
		condition += " AND stream_path=?"
		values = append(values, req.StreamPath)
	}
	if req.RecordMode != "" {
		condition += " AND record_mode=?"
		values = append(values, req.RecordMode)
	}
	values = append([]any{condition}, values...)
	err = countQuery.Find(&streams, values...).Count(&totalCount).Error
	if err != nil {
		return
	}
	query.Find(&streams, values...)
	resp = &pb.ResponseList{
		PageSize:   req.PageSize,
		PageNum:    req.PageNum,
		TotalCount: uint32(totalCount),
	}
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
	if p.DB == nil {
		err = pkg.ErrNoDB
		return
	}
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
	if p.DB == nil {
		err = pkg.ErrNoDB
		return
	}
	ids := req.GetIds()
	var result []*m7s.RecordStream
	if len(ids) > 0 {
		p.DB.Find(&result, "stream_path=? AND id IN ?", req.StreamPath, ids)
	} else {
		startTime, endTime, err := util.TimeRangeQueryParse(url.Values{"range": []string{req.Range}, "start": []string{req.StartTime}, "end": []string{req.EndTime}})
		if err != nil {
			return nil, err
		}
		p.DB.Find(&result, "stream_path=? AND start_time>=? AND end_time<=?", req.StreamPath, startTime, endTime)
	}
	err = p.DB.Delete(result).Error
	if err != nil {
		return
	}
	var apiResult []*pb.RecordFile
	for _, recordFile := range result {
		apiResult = append(apiResult, &pb.RecordFile{
			Id:         uint32(recordFile.ID),
			StartTime:  timestamppb.New(recordFile.StartTime),
			EndTime:    timestamppb.New(recordFile.EndTime),
			FilePath:   recordFile.FilePath,
			StreamPath: recordFile.StreamPath,
		})
		err = os.Remove(recordFile.FilePath)
		if err != nil {
			return
		}
	}
	resp = &pb.ResponseDelete{
		Data: apiResult,
	}
	return
}

func (p *MP4Plugin) download(w http.ResponseWriter, r *http.Request) {
	if p.DB == nil {
		http.Error(w, pkg.ErrNoDB.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	streamPath := r.PathValue("streamPath")

	query := r.URL.Query()
	var streams []m7s.RecordStream
	if id := query.Get("id"); id != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s.mp4", streamPath, id))
		p.DB.Find(&streams, "id=? AND stream_path=?", id, streamPath)
		if len(streams) == 0 {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, streams[0].FilePath)
		return
	}

	startTime, endTime, err := util.TimeRangeQueryParse(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.Info("download", "streamPath", streamPath, "start", startTime, "end", endTime)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_%s.mp4", streamPath, startTime.Format("20060102150405"), endTime.Format("20060102150405")))

	p.DB.Find(&streams, "end_time>? AND start_time<? AND stream_path=? AND record_mode=0", startTime, endTime, streamPath)
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

func (p *MP4Plugin) EventStart(ctx context.Context, req *pb.ReqEventRecord) (res *pb.ResponseEventRecord, err error) {
	beforeDuration := p.BeforeDuration
	afterDuration := p.AfterDuration

	if req.BeforeDuration != "" {
		beforeDuration, err = time.ParseDuration(req.BeforeDuration)
		if err != nil {
			p.Info("error", err)
		}
	}
	if req.AfterDuration != "" {
		afterDuration, err = time.ParseDuration(req.AfterDuration)
		if err != nil {
			p.Info("error", err)
		}
	}
	recorder := p.Meta.Recorder()
	recorderJobs := p.Server.Records
	var tmpJob *m7s.RecordJob
	if recorderJobs.Length > 0 {
		for job := range recorderJobs.Range {
			if job.StreamPath == req.StreamPath {
				tmpJob = job
			}
		}
	}
	if tmpJob == nil { //为空表示没有正在进行的录制，也就是没有自动录像，则进行正常的事件录像
		if stream, ok := p.Server.Streams.Get(req.StreamPath); ok {
			recordConf := config.Record{
				Append:   false,
				Fragment: 0,
				FilePath: filepath.Join(p.EventRecordFilePath, stream.StreamPath, time.Now().Local().Format("2006-01-02-15-04-05")),
			}
			recordJob := recorder.GetRecordJob()
			recordJob.EventId = req.EventId
			recordJob.EventLevel = req.EventLevel
			recordJob.EventName = req.EventName
			recordJob.EventDesc = req.EventDesc
			recordJob.AfterDuration = afterDuration
			recordJob.BeforeDuration = beforeDuration
			recordJob.RecordMode = "1"
			var subconfig config.Subscribe
			defaults.SetDefaults(&subconfig)
			subconfig.BufferTime = beforeDuration
			job := recorder.GetRecordJob().Init(recorder, &p.Plugin, stream.StreamPath, recordConf, &subconfig)
			job.Depend(stream)
		}
	} else {
		if tmpJob.AfterDuration != 0 { //当前有事件录像正在录制，则更新该录像的结束时间
			tmpJob.AfterDuration = time.Duration(tmpJob.Subscriber.VideoReader.AbsTime)*time.Millisecond + afterDuration
		} else { //当前有自动录像正在录制，则生成事件录像的记录，而不去生成事件录像的文件
			recordStream := &m7s.RecordStream{
				StreamPath:     req.StreamPath,
				EventId:        req.EventId,
				EventLevel:     req.EventLevel,
				EventDesc:      req.EventDesc,
				EventName:      req.EventName,
				RecordMode:     "1",
				BeforeDuration: beforeDuration,
				AfterDuration:  afterDuration,
			}
			now := time.Now()
			startTime := now.Add(-beforeDuration)
			endTime := now.Add(afterDuration)
			recordStream.StartTime = startTime
			recordStream.EndTime = endTime
			if p.DB != nil {
				p.DB.Save(&recordStream)
			}
		}
	}
	return res, err
}
