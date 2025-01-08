package flv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"gorm.io/gorm"
	"m7s.live/v5"
	"m7s.live/v5/pkg"
	"m7s.live/v5/pkg/task"
	"m7s.live/v5/pkg/util"
	rtmp "m7s.live/v5/plugin/rtmp/pkg"
)

type WriteFlvMetaTagQueueTask struct {
	task.Work
}

var writeMetaTagQueueTask WriteFlvMetaTagQueueTask

func init() {
	m7s.Servers.AddTask(&writeMetaTagQueueTask)
}

type writeMetaTagTask struct {
	task.Task
	file     *os.File
	writer   *FlvWriter
	flags    byte
	metaData []byte
}

func (task *writeMetaTagTask) Start() (err error) {
	defer func() {
		err = task.file.Close()
		if info, err := task.file.Stat(); err == nil && info.Size() == 0 {
			err = os.Remove(info.Name())
		}
	}()
	var tempFile *os.File
	if tempFile, err = os.CreateTemp("", "*.flv"); err != nil {
		task.Error("create temp file failed", "err", err)
		return
	} else {
		defer func() {
			err = tempFile.Close()
			err = os.Remove(tempFile.Name())
			task.Info("writeMetaData success")
		}()
		_, err = tempFile.Write([]byte{'F', 'L', 'V', 0x01, task.flags, 0, 0, 0, 9, 0, 0, 0, 0})
		if err != nil {
			task.Error(err.Error())
			return
		}
		task.writer = NewFlvWriter(tempFile)
		err = task.writer.WriteTag(FLV_TAG_TYPE_SCRIPT, 0, uint32(len(task.metaData)), task.metaData)
		_, err = task.file.Seek(13, io.SeekStart)
		if err != nil {
			task.Error("writeMetaData Seek failed", "err", err)
			return
		}
		_, err = io.Copy(tempFile, task.file)
		if err != nil {
			task.Error("writeMetaData Copy failed", "err", err)
			return
		}
		_, err = tempFile.Seek(0, io.SeekStart)
		_, err = task.file.Seek(0, io.SeekStart)
		_, err = io.Copy(task.file, tempFile)
		if err != nil {
			task.Error("writeMetaData Copy failed", "err", err)
		}
		return
	}
}

func writeMetaTag(file *os.File, suber *m7s.Subscriber, filepositions []uint64, times []float64, duration *int64) {
	ar, vr := suber.AudioReader, suber.VideoReader
	hasAudio, hasVideo := ar != nil, vr != nil
	var amf rtmp.AMF
	metaData := rtmp.EcmaArray{
		"MetaDataCreator": "m7s/" + m7s.Version,
		"hasVideo":        hasVideo,
		"hasAudio":        hasAudio,
		"hasMatadata":     true,
		"canSeekToEnd":    true,
		"duration":        float64(*duration) / 1000,
		"hasKeyFrames":    len(filepositions) > 0,
		"filesize":        0,
	}
	var flags byte
	if hasAudio {
		ctx := ar.Track.ICodecCtx.GetBase().(pkg.IAudioCodecCtx)
		flags |= (1 << 2)
		metaData["audiocodecid"] = int(rtmp.ParseAudioCodec(ctx.FourCC()))
		metaData["audiosamplerate"] = ctx.GetSampleRate()
		metaData["audiosamplesize"] = ctx.GetSampleSize()
		metaData["stereo"] = ctx.GetChannels() == 2
	}
	if hasVideo {
		ctx := vr.Track.ICodecCtx.GetBase().(pkg.IVideoCodecCtx)
		flags |= 1
		metaData["videocodecid"] = int(rtmp.ParseVideoCodec(ctx.FourCC()))
		metaData["width"] = ctx.Width()
		metaData["height"] = ctx.Height()
		metaData["framerate"] = vr.Track.FPS
		metaData["videodatarate"] = vr.Track.BPS
		metaData["keyframes"] = map[string]any{
			"filepositions": filepositions,
			"times":         times,
		}
	}
	amf.Marshals("onMetaData", metaData)
	offset := amf.Len() + 13 + 15
	if keyframesCount := len(filepositions); keyframesCount > 0 {
		metaData["filesize"] = uint64(offset) + filepositions[keyframesCount-1]
		for i := range filepositions {
			filepositions[i] += uint64(offset)
		}
		metaData["keyframes"] = map[string]any{
			"filepositions": filepositions,
			"times":         times,
		}
	}
	amf.Reset()
	marshals := amf.Marshals("onMetaData", metaData)
	task := &writeMetaTagTask{
		file:     file,
		flags:    flags,
		metaData: marshals,
	}
	task.Logger = suber.Logger.With("file", file.Name())
	writeMetaTagQueueTask.AddTask(task)
}

func NewRecorder() m7s.IRecorder {
	return &Recorder{}
}

type Recorder struct {
	m7s.DefaultRecorder
	stream m7s.RecordStream
}

var CustomFileName = func(job *m7s.RecordJob) string {
	if job.Fragment == 0 || job.Append {
		return fmt.Sprintf("%s.flv", job.FilePath)
	}
	return filepath.Join(job.FilePath, fmt.Sprintf("%d.flv", time.Now().Unix()))
}

func (r *Recorder) createStream(start time.Time) (err error) {
	recordJob := &r.RecordJob
	sub := recordJob.Subscriber
	r.stream = m7s.RecordStream{
		StartTime:      start,
		StreamPath:     sub.StreamPath,
		FilePath:       CustomFileName(&r.RecordJob),
		EventId:        recordJob.EventId,
		EventDesc:      recordJob.EventDesc,
		EventName:      recordJob.EventName,
		EventLevel:     recordJob.EventLevel,
		BeforeDuration: recordJob.BeforeDuration,
		AfterDuration:  recordJob.AfterDuration,
		Mode:           recordJob.Mode,
		Type:           "flv",
	}
	dir := filepath.Dir(r.stream.FilePath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	if sub.Publisher.HasAudioTrack() {
		r.stream.AudioCodec = sub.Publisher.AudioTrack.ICodecCtx.FourCC().String()
	}
	if sub.Publisher.HasVideoTrack() {
		r.stream.VideoCodec = sub.Publisher.VideoTrack.ICodecCtx.FourCC().String()
	}
	if recordJob.Plugin.DB != nil {
		recordJob.Plugin.DB.Save(&r.stream)
	}
	return
}

func (r *Recorder) writeTailer(end time.Time) {
	if r.stream.EndTime.After(r.stream.StartTime) {
		return
	}
	r.stream.EndTime = end
	if r.RecordJob.Plugin.DB != nil {
		r.RecordJob.Plugin.DB.Save(&r.stream)
		writeMetaTagQueueTask.AddTask(&eventRecordCheck{
			DB:         r.RecordJob.Plugin.DB,
			streamPath: r.stream.StreamPath,
		})
	}
}

func (r *Recorder) Dispose() {
	r.writeTailer(time.Now())
}

type eventRecordCheck struct {
	task.Task
	DB         *gorm.DB
	streamPath string
}

func (t *eventRecordCheck) Run() (err error) {
	var eventRecordStreams []m7s.RecordStream
	queryRecord := m7s.RecordStream{
		EventLevel: m7s.EventLevelHigh,
		Mode:       m7s.RecordModeEvent,
		Type:       "flv",
	}
	t.DB.Where(&queryRecord).Find(&eventRecordStreams, "stream_path=?", t.streamPath) //搜索事件录像，且为重要事件（无法自动删除）
	if len(eventRecordStreams) > 0 {
		for _, recordStream := range eventRecordStreams {
			var unimportantEventRecordStreams []m7s.RecordStream
			queryRecord.EventLevel = m7s.EventLevelLow
			query := `(start_time BETWEEN ? AND ?)
							OR (end_time BETWEEN ? AND ?) 
							OR (? BETWEEN start_time AND end_time) 
							OR (? BETWEEN start_time AND end_time) AND stream_path=? `
			t.DB.Where(&queryRecord).Where(query, recordStream.StartTime, recordStream.EndTime, recordStream.StartTime, recordStream.EndTime, recordStream.StartTime, recordStream.EndTime, recordStream.StreamPath).Find(&unimportantEventRecordStreams)
			if len(unimportantEventRecordStreams) > 0 {
				for _, unimportantEventRecordStream := range unimportantEventRecordStreams {
					unimportantEventRecordStream.EventLevel = m7s.EventLevelHigh
					t.DB.Save(&unimportantEventRecordStream)
				}
			}
		}
	}
	return
}

func (r *Recorder) Run() (err error) {
	var file *os.File
	var filepositions []uint64
	var times []float64
	var offset int64
	var duration int64
	ctx := &r.RecordJob
	suber := ctx.Subscriber
	noFragment := ctx.Fragment == 0 || ctx.Append
	startTime := time.Now()
	if ctx.BeforeDuration > 0 {
		startTime = startTime.Add(-ctx.BeforeDuration)
	}
	if err = r.createStream(startTime); err != nil {
		return
	}
	if noFragment {
		file, err = os.OpenFile(r.stream.FilePath, os.O_CREATE|os.O_RDWR|util.Conditional(ctx.Append, os.O_APPEND, os.O_TRUNC), 0666)
		if err != nil {
			return
		}
		defer writeMetaTag(file, suber, filepositions, times, &duration)
	}
	if ctx.Append {
		var metaData rtmp.EcmaArray
		metaData, err = ReadMetaData(file)
		keyframes := metaData["keyframes"].(map[string]any)
		filepositions = slices.Collect(func(yield func(uint64) bool) {
			for _, v := range keyframes["filepositions"].([]float64) {
				yield(uint64(v))
			}
		})
		times = keyframes["times"].([]float64)
		if _, err = file.Seek(-4, io.SeekEnd); err != nil {
			ctx.Error("seek file failed", "err", err)
			_, err = file.Write(FLVHead)
		} else {
			tmp := make(util.Buffer, 4)
			tmp2 := tmp
			_, err = file.Read(tmp)
			tagSize := tmp.ReadUint32()
			tmp = tmp2
			_, err = file.Seek(int64(tagSize), io.SeekEnd)
			_, err = file.Read(tmp2)
			ts := tmp2.ReadUint24() | (uint32(tmp[3]) << 24)
			ctx.Info("append flv", "last tagSize", tagSize, "last ts", ts)
			suber.StartAudioTS = time.Duration(ts) * time.Millisecond
			suber.StartVideoTS = time.Duration(ts) * time.Millisecond
			offset, err = file.Seek(0, io.SeekEnd)
		}
	} else if ctx.Fragment == 0 {
		_, err = file.Write(FLVHead)
	} else {
		if file, err = os.OpenFile(r.stream.FilePath, os.O_CREATE|os.O_RDWR, 0666); err != nil {
			return
		}
		_, err = file.Write(FLVHead)
	}
	writer := NewFlvWriter(file)
	checkFragment := func(absTime uint32) {
		if duration = int64(absTime); time.Duration(duration)*time.Millisecond >= ctx.Fragment {
			writeMetaTag(file, suber, filepositions, times, &duration)
			r.writeTailer(time.Now())
			filepositions = []uint64{0}
			times = []float64{0}
			offset = 0
			if err = r.createStream(time.Now()); err != nil {
				return
			}
			if file, err = os.OpenFile(r.stream.FilePath, os.O_CREATE|os.O_RDWR, 0666); err != nil {
				return
			}
			_, err = file.Write(FLVHead)
			writer = NewFlvWriter(file)
			if vr := suber.VideoReader; vr != nil {
				vr.ResetAbsTime()
				seq := vr.Track.SequenceFrame.(*rtmp.RTMPVideo)
				err = writer.WriteTag(FLV_TAG_TYPE_VIDEO, 0, uint32(seq.Size), seq.Buffers...)
				offset = int64(seq.Size + 15)
			}
			if ar := suber.AudioReader; ar != nil {
				ar.ResetAbsTime()
				if ar.Track.SequenceFrame != nil {
					seq := ar.Track.SequenceFrame.(*rtmp.RTMPAudio)
					err = writer.WriteTag(FLV_TAG_TYPE_AUDIO, 0, uint32(seq.Size), seq.Buffers...)
					offset += int64(seq.Size + 15)
				}
			}
		}
	}

	return m7s.PlayBlock(ctx.Subscriber, func(audio *rtmp.RTMPAudio) (err error) {
		if suber.VideoReader == nil && !noFragment {
			checkFragment(suber.AudioReader.AbsTime)
		}
		err = writer.WriteTag(FLV_TAG_TYPE_AUDIO, suber.AudioReader.AbsTime, uint32(audio.Size), audio.Buffers...)
		offset += int64(audio.Size + 15)
		return
	}, func(video *rtmp.RTMPVideo) (err error) {
		if suber.VideoReader.Value.IDR {
			filepositions = append(filepositions, uint64(offset))
			times = append(times, float64(suber.VideoReader.AbsTime)/1000)
			if !noFragment {
				checkFragment(suber.VideoReader.AbsTime)
			}
		}
		err = writer.WriteTag(FLV_TAG_TYPE_VIDEO, suber.VideoReader.AbsTime, uint32(video.Size), video.Buffers...)
		offset += int64(video.Size + 15)
		return
	})
}
