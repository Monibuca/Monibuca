package m7s

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sync"
	"time"

	"m7s.live/v5/pkg"
	"m7s.live/v5/pkg/codec"
	"m7s.live/v5/pkg/task"

	. "m7s.live/v5/pkg"
	"m7s.live/v5/pkg/config"
	"m7s.live/v5/pkg/util"
)

type PublisherState int

const (
	PublisherStateInit PublisherState = iota
	PublisherStateTrackAdded
	PublisherStateSubscribed
	PublisherStateWaitSubscriber
	PublisherStateDisposed
)

const (
	PublishTypePull      = "pull"
	PublishTypeServer    = "server"
	PublishTypeVod       = "vod"
	PublishTypeTransform = "transform"
	PublishTypeReplay    = "replay"
)

const threshold = 10 * time.Millisecond

type SpeedControl struct {
	speed          float64
	pausedTime     time.Duration
	beginTime      time.Time
	beginTimestamp time.Duration
	Delta          time.Duration
}

func (s *SpeedControl) speedControl(speed float64, ts time.Duration) {
	if speed != s.speed || s.beginTime.IsZero() {
		s.speed = speed
		s.beginTime = time.Now()
		s.beginTimestamp = ts
		s.pausedTime = 0
	} else {
		elapsed := time.Since(s.beginTime) - s.pausedTime
		if speed == 0 {
			s.Delta = ts - elapsed
			return
		}
		should := time.Duration(float64(ts-s.beginTimestamp) / speed)
		s.Delta = should - elapsed
		//fmt.Println(speed, elapsed, should, s.Delta)
		if s.Delta > threshold {
			time.Sleep(min(s.Delta, time.Millisecond*500))
		}
	}
}

type AVTracks struct {
	*AVTrack
	SpeedControl
	util.Collection[reflect.Type, *AVTrack]
	sync.RWMutex
	baseTs time.Duration //from old publisher's lastTs
}

func (t *AVTracks) Set(track *AVTrack) {
	t.Lock()
	defer t.Unlock()
	t.AVTrack = track
	track.BaseTs = t.baseTs
	t.Add(track)
}

func (t *AVTracks) SetMinBuffer(start time.Duration) {
	if t.AVTrack == nil {
		return
	}
	t.AVTrack.BufferRange[0] = start
}

func (t *AVTracks) GetOrCreate(dataType reflect.Type) *AVTrack {
	t.Lock()
	defer t.Unlock()
	if track, ok := t.Get(dataType); ok {
		return track
	}
	if t.AVTrack == nil {
		return nil
	}
	return t.CreateSubTrack(dataType)
}

func (t *AVTracks) CheckTimeout(timeout time.Duration) bool {
	if t.AVTrack == nil {
		return false
	}
	return time.Since(t.AVTrack.LastValue.WriteTime) > timeout
}

func (t *AVTracks) CreateSubTrack(dataType reflect.Type) (track *AVTrack) {
	track = NewAVTrack(dataType, t.AVTrack)
	track.WrapIndex = t.Length
	t.Add(track)
	return
}

func (t *AVTracks) Dispose() {
	t.Lock()
	defer t.Unlock()
	for track := range t.Range {
		track.Ready(ErrDiscard)
		if track == t.AVTrack || track.RingWriter != t.AVTrack.RingWriter {
			track.Dispose()
		}
	}
	t.AVTrack = nil
	t.Clear()
}

type Publisher struct {
	PubSubBase
	config.Publish
	State                  PublisherState
	Paused                 *util.Promise
	pauseTime              time.Time
	AudioTrack, VideoTrack AVTracks
	audioReady, videoReady *util.Promise
	TimeoutTimer           *time.Timer
	DataTrack              *DataTrack
	Subscribers            SubscriberCollection
	GOP                    int
	OnSeek                 func(time.Time)
	OnGetPosition          func() time.Time
	PullProxy              *PullProxy
	dumpFile               *os.File
	MaxFPS                 float64
	dropRate               float64 // 丢帧率，0-1之间
}

func (p *Publisher) SubscriberRange(yield func(sub *Subscriber) bool) {
	p.Subscribers.Range(yield)
}

func (p *Publisher) GetKey() string {
	return p.StreamPath
}

// createPublisher -> Start -> WriteAudio/WriteVideo -> Dispose
func createPublisher(p *Plugin, streamPath string, conf config.Publish) (publisher *Publisher) {
	publisher = &Publisher{Publish: conf}
	publisher.Type = conf.PubType
	publisher.ID = task.GetNextTaskID()
	publisher.Plugin = p
	publisher.TimeoutTimer = time.NewTimer(p.config.PublishTimeout)
	publisher.Logger = p.Logger.With("streamPath", streamPath, "pId", publisher.ID)
	publisher.Init(streamPath, &publisher.Publish)
	return
}

func (p *Publisher) Start() (err error) {
	s := p.Plugin.Server
	if oldPublisher, ok := s.Streams.Get(p.StreamPath); ok {
		if p.KickExist {
			p.takeOver(oldPublisher)
		} else {
			return ErrStreamExist
		}
	}
	s.Streams.Set(p)
	p.Info("publish")
	p.processPullProxyOnStart()
	p.audioReady = util.NewPromiseWithTimeout(p, p.PublishTimeout)
	if !p.PubAudio {
		p.audioReady.Reject(ErrMuted)
	}
	p.videoReady = util.NewPromiseWithTimeout(p, p.PublishTimeout)
	if !p.PubVideo {
		p.videoReady.Reject(ErrMuted)
	}
	if p.Dump {
		f := filepath.Join("./dump", p.StreamPath)
		os.MkdirAll(filepath.Dir(f), 0666)
		p.dumpFile, _ = os.OpenFile(filepath.Join("./dump", p.StreamPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	}

	s.Waiting.WakeUp(p.StreamPath, p)
	p.processAliasOnStart()
	for plugin := range s.Plugins.Range {
		plugin.OnPublish(p)
	}
	//s.Transforms.PublishEvent <- p
	p.AddTask(&PublishTimeout{Publisher: p})
	if p.PublishTimeout > 0 {
		p.AddTask(&PublishNoDataTimeout{Publisher: p})
	}
	return
}

type PublishTimeout struct {
	task.ChannelTask
	Publisher *Publisher
}

func (p *PublishTimeout) Start() error {
	p.SignalChan = p.Publisher.TimeoutTimer.C
	return nil
}

func (p *PublishTimeout) Dispose() {
	p.Publisher.TimeoutTimer.Stop()
}

func (p *PublishTimeout) Tick(any) {
	if p.Publisher.Paused != nil {
		return
	}
	switch p.Publisher.State {
	case PublisherStateInit:
		if p.Publisher.PublishTimeout > 0 {
			p.Publisher.Stop(ErrPublishTimeout)
		}
	case PublisherStateTrackAdded:
		if p.Publisher.Publish.IdleTimeout > 0 {
			p.Publisher.Stop(ErrPublishIdleTimeout)
		}
	case PublisherStateSubscribed:
	case PublisherStateWaitSubscriber:
		if p.Publisher.Publish.DelayCloseTimeout > 0 {
			p.Publisher.Stop(ErrPublishDelayCloseTimeout)
		}
	}
}

type PublishNoDataTimeout struct {
	task.TickTask
	Publisher *Publisher
}

func (p *PublishNoDataTimeout) GetTickInterval() time.Duration {
	return time.Second * 5
}

func (p *PublishNoDataTimeout) Tick(any) {
	if p.Publisher.Paused != nil {
		return
	}
	if p.Publisher.VideoTrack.CheckTimeout(p.Publisher.PublishTimeout) {
		p.Error("video timeout", "writeTime", p.Publisher.VideoTrack.LastValue.WriteTime)
		p.Publisher.Stop(ErrPublishTimeout)
	}
	if p.Publisher.AudioTrack.CheckTimeout(p.Publisher.PublishTimeout) {
		p.Error("audio timeout", "writeTime", p.Publisher.AudioTrack.LastValue.WriteTime)
		p.Publisher.Stop(ErrPublishTimeout)
	}
}

func (p *Publisher) RemoveSubscriber(subscriber *Subscriber) {
	p.Subscribers.Remove(subscriber)
	p.Info("subscriber -1", "count", p.Subscribers.Length)
	if p.Plugin == nil {
		return
	}
	if subscriber.BufferTime == p.BufferTime && p.Subscribers.Length > 0 {
		p.BufferTime = slices.MaxFunc(p.Subscribers.Items, func(a, b *Subscriber) int {
			return int(a.BufferTime - b.BufferTime)
		}).BufferTime
	} else {
		p.BufferTime = p.Plugin.GetCommonConf().Publish.BufferTime
	}
	p.AudioTrack.SetMinBuffer(p.BufferTime)
	p.VideoTrack.SetMinBuffer(p.BufferTime)
	if p.State == PublisherStateSubscribed && p.Subscribers.Length == 0 {
		p.State = PublisherStateWaitSubscriber
		if p.DelayCloseTimeout > 0 {
			p.TimeoutTimer.Reset(p.DelayCloseTimeout)
		}
	}
}

func (p *Publisher) AddSubscriber(subscriber *Subscriber) {
	oldPublisher := subscriber.Publisher
	subscriber.Publisher = p
	if oldPublisher == nil {
		close(subscriber.waitPublishDone)
	} else {
		if subscriber.waitingPublish() {
			subscriber.Info("publisher recover", "pid", p.ID)
		} else {
			subscriber.Info("publisher changed", "prePid", oldPublisher.ID, "pid", p.ID)
		}
	}
	subscriber.waitStartTime = time.Time{}
	if p.Subscribers.AddUnique(subscriber) {
		p.Info("subscriber +1", "count", p.Subscribers.Length)
		if subscriber.BufferTime > p.BufferTime {
			p.BufferTime = subscriber.BufferTime
			p.AudioTrack.SetMinBuffer(p.BufferTime)
			p.VideoTrack.SetMinBuffer(p.BufferTime)
		}
		switch p.State {
		case PublisherStateTrackAdded, PublisherStateWaitSubscriber:
			p.State = PublisherStateSubscribed
			if p.PublishTimeout > 0 {
				p.TimeoutTimer.Reset(p.PublishTimeout)
			}
		}
	}
}

func (p *Publisher) writeAV(t *AVTrack, data IAVFrame) {
	frame := &t.Value
	frame.Wraps = append(frame.Wraps, data)
	ts := data.GetTimestamp()
	frame.CTS = data.GetCTS()
	bytesIn := frame.Wraps[0].GetSize()
	t.AddBytesIn(bytesIn)
	frame.Timestamp = t.Tame(ts, t.FPS, p.Scale)
	if p.Enabled(p, task.TraceLevel) {
		codec := t.FourCC().String()
		data := frame.Wraps[0].String()
		p.Trace("write", "seq", frame.Sequence, "baseTs", int32(t.BaseTs/time.Millisecond), "ts0", uint32(ts/time.Millisecond), "ts", uint32(frame.Timestamp/time.Millisecond), "codec", codec, "size", bytesIn, "data", data)
	}
}

func (p *Publisher) trackAdded() error {
	if p.Subscribers.Length > 0 {
		p.State = PublisherStateSubscribed
	} else {
		p.State = PublisherStateTrackAdded
	}
	return nil
}

func (p *Publisher) dropFrame(t *AVTrack, idr *util.Ring[AVFrame]) (drop bool) {
	// Frame dropping logic based on MaxFPS
	if p.MaxFPS > 0 && float64(t.FPS) > p.MaxFPS {
		dropRatio := float64(t.FPS)/p.MaxFPS - 1 // How many frames to drop for each frame kept
		p.dropRate = dropRatio / (dropRatio + 1) // 计算丢帧率
		if p.Scale >= 8 {
			// Only keep I-frames when Scale >= 8
			if !t.Value.IDR {
				// Drop all P-frames
				return true
			}
			// Drop I-frames based on FPS ratio
			if dropRatio > 1 && t.Value.IDR && t.Value.Sequence%uint32(dropRatio+1) != 0 {
				return true
			}
		} else {
			// Normal frame dropping strategy
			if !t.Value.IDR {
				// Drop P-frames based on position in GOP and FPS ratio
				if idr != nil {
					posInGOP := int(t.Value.Sequence - idr.Value.Sequence)
					gopThreshold := int(float64(p.GOP) / (dropRatio + 1))
					if posInGOP > gopThreshold {
						return true
					}
				}
			}
		}
	} else {
		p.dropRate = 0
	}
	return
}

func (p *Publisher) WriteVideo(data IAVFrame) (err error) {
	defer func() {
		if err != nil {
			data.Recycle()
		}
	}()
	if err = p.Err(); err != nil {
		return
	}
	if p.dumpFile != nil {
		data.Dump(1, p.dumpFile)
	}
	if !p.PubVideo {
		return ErrMuted
	}
	t := p.VideoTrack.AVTrack
	if t == nil {
		t = NewAVTrack(data, p.Logger.With("track", "video"), &p.Publish, p.videoReady)
		p.VideoTrack.Set(t)
		p.Call(p.trackAdded)
	}
	oldCodecCtx := t.ICodecCtx
	err = data.Parse(t)
	codecCtxChanged := oldCodecCtx != t.ICodecCtx
	if err != nil {
		p.Error("parse", "err", err)
		return err
	}

	if t.ICodecCtx == nil {
		return ErrUnsupportCodec
	}
	if codecCtxChanged && oldCodecCtx != nil {
		oldWidth, oldHeight := oldCodecCtx.(IVideoCodecCtx).Width(), oldCodecCtx.(IVideoCodecCtx).Height()
		newWidth, newHeight := t.ICodecCtx.(IVideoCodecCtx).Width(), t.ICodecCtx.(IVideoCodecCtx).Height()
		if oldWidth != newWidth || oldHeight != newHeight {
			p.Info("video resolution changed", "oldWidth", oldWidth, "oldHeight", oldHeight, "newWidth", newWidth, "newHeight", newHeight)
		}
	}
	var idr *util.Ring[AVFrame]
	if t.IDRingList.Len() > 0 {
		idr = t.IDRingList.Back().Value
		if p.dropFrame(t, idr) {
			data.Recycle()
			return nil
		}
	}
	if t.Value.IDR {
		if !t.IsReady() {
			t.Ready(nil)
		} else if idr != nil {
			p.GOP = int(t.Value.Sequence - idr.Value.Sequence)
		} else {
			p.GOP = 0
		}
		if p.AudioTrack.Length > 0 {
			p.AudioTrack.PushIDR()
		}
	}
	p.writeAV(t, data)
	if p.VideoTrack.Length > 1 && p.VideoTrack.IsReady() {
		if t.Value.Raw == nil {
			if err = t.Value.Demux(t.ICodecCtx); err != nil {
				t.Error("to raw", "err", err)
				return err
			}
		}
		for i, track := range p.VideoTrack.Items[1:] {
			toType := track.FrameType.Elem()
			toFrame := reflect.New(toType).Interface().(IAVFrame)
			if track.ICodecCtx == nil {
				if track.ICodecCtx, track.SequenceFrame, err = toFrame.ConvertCtx(t.ICodecCtx); err != nil {
					track.Error("DecodeConfig", "err", err)
					return
				}
				if t.IDRingList.Len() > 0 {
					for rf := t.IDRingList.Front().Value; rf != t.Ring; rf = rf.Next() {
						if i == 0 && rf.Value.Raw == nil {
							if err = rf.Value.Demux(t.ICodecCtx); err != nil {
								t.Error("to raw", "err", err)
								return err
							}
						}
						toFrame := reflect.New(toType).Interface().(IAVFrame)
						toFrame.SetAllocator(data.GetAllocator())
						toFrame.Mux(track.ICodecCtx, &rf.Value)
						rf.Value.Wraps = append(rf.Value.Wraps, toFrame)
					}
				}
			}
			toFrame.SetAllocator(data.GetAllocator())
			toFrame.Mux(track.ICodecCtx, &t.Value)
			if codecCtxChanged {
				track.ICodecCtx, track.SequenceFrame, err = toFrame.ConvertCtx(t.ICodecCtx)
			}
			t.Value.Wraps = append(t.Value.Wraps, toFrame)
			if track.ICodecCtx != nil {
				track.Ready(err)
			}
		}
	}
	t.Step()

	p.VideoTrack.speedControl(p.Speed, t.LastTs)
	return
}

func (p *Publisher) WriteAudio(data IAVFrame) (err error) {
	defer func() {
		if err != nil {
			data.Recycle()
		}
	}()
	if err = p.Err(); err != nil {
		return
	}
	if p.dumpFile != nil {
		data.Dump(0, p.dumpFile)
	}
	if !p.PubAudio {
		return ErrMuted
	}
	// 根据丢帧率进行音频帧丢弃
	if p.dropRate > 0 {
		t := p.AudioTrack.AVTrack
		if t != nil {
			// 使用序列号进行平均丢帧
			if t.Value.Sequence%uint32(1/p.dropRate) != 0 {
				data.Recycle()
				return nil
			}
		}
	}
	t := p.AudioTrack.AVTrack
	if t == nil {
		t = NewAVTrack(data, p.Logger.With("track", "audio"), &p.Publish, p.audioReady)
		p.AudioTrack.Set(t)
		p.Call(p.trackAdded)
	}
	oldCodecCtx := t.ICodecCtx
	err = data.Parse(t)
	codecCtxChanged := oldCodecCtx != t.ICodecCtx
	if t.ICodecCtx == nil {
		return ErrUnsupportCodec
	}
	t.Ready(err)
	p.writeAV(t, data)
	if p.AudioTrack.Length > 1 && p.AudioTrack.IsReady() {
		if t.Value.Raw == nil {
			if err = t.Value.Demux(t.ICodecCtx); err != nil {
				t.Error("to raw", "err", err)
				return err
			}
		}
		for i, track := range p.AudioTrack.Items[1:] {
			toType := track.FrameType.Elem()
			toFrame := reflect.New(toType).Interface().(IAVFrame)
			if track.ICodecCtx == nil {
				if track.ICodecCtx, track.SequenceFrame, err = toFrame.ConvertCtx(t.ICodecCtx); err != nil {
					track.Error("DecodeConfig", "err", err)
					return
				}
				if idr := p.AudioTrack.GetOldestIDR(); idr != nil {
					for rf := idr; rf != t.Ring; rf = rf.Next() {
						if i == 0 && rf.Value.Raw == nil {
							if err = rf.Value.Demux(t.ICodecCtx); err != nil {
								t.Error("to raw", "err", err)
								return err
							}
						}
						toFrame := reflect.New(toType).Interface().(IAVFrame)
						toFrame.SetAllocator(data.GetAllocator())
						toFrame.Mux(track.ICodecCtx, &rf.Value)
						rf.Value.Wraps = append(rf.Value.Wraps, toFrame)
					}
				}
			}
			toFrame.SetAllocator(data.GetAllocator())
			toFrame.Mux(track.ICodecCtx, &t.Value)
			if codecCtxChanged {
				track.ICodecCtx, track.SequenceFrame, err = toFrame.ConvertCtx(t.ICodecCtx)
			}
			t.Value.Wraps = append(t.Value.Wraps, toFrame)
			if track.ICodecCtx != nil {
				track.Ready(err)
			}
		}
	}
	t.Step()
	p.AudioTrack.speedControl(p.Speed, t.LastTs)
	return
}

func (p *Publisher) WriteData(data IDataFrame) (err error) {
	for subscriber := range p.SubscriberRange {
		if subscriber.DataChannel == nil {
			continue
		}
		select {
		case subscriber.DataChannel <- data:
		default:
			p.Warn("subscriber channel full", "subscriber", subscriber.ID)
		}
	}
	return nil
}

func (p *Publisher) GetAudioCodecCtx() (ctx codec.ICodecCtx) {
	if p.HasAudioTrack() {
		return p.AudioTrack.ICodecCtx
	}
	return nil
}

func (p *Publisher) GetVideoCodecCtx() (ctx codec.ICodecCtx) {
	if p.HasVideoTrack() {
		return p.VideoTrack.ICodecCtx
	}
	return nil
}

func (p *Publisher) GetAudioTrack(dataType reflect.Type) (t *AVTrack) {
	return p.AudioTrack.GetOrCreate(dataType)
}

func (p *Publisher) GetVideoTrack(dataType reflect.Type) (t *AVTrack) {
	return p.VideoTrack.GetOrCreate(dataType)
}

func (p *Publisher) HasAudioTrack() bool {
	return p.AudioTrack.Length > 0
}

func (p *Publisher) HasVideoTrack() bool {
	return p.VideoTrack.Length > 0
}

func (p *Publisher) Dispose() {
	s := p.Plugin.Server
	if !p.StopReasonIs(ErrKick) {
		s.Streams.Remove(p)
	}
	if p.Paused != nil {
		p.Paused.Reject(p.StopReason())
	}
	p.processAliasOnDispose()
	p.AudioTrack.Dispose()
	p.VideoTrack.Dispose()
	p.Info("unpublish", "remain", s.Streams.Length, "reason", p.StopReason())
	if p.dumpFile != nil {
		p.dumpFile.Close()
	}
	p.State = PublisherStateDisposed
	p.processPullProxyOnDispose()
}

func (p *Publisher) TransferSubscribers(newPublisher *Publisher) {
	p.Info("transfer subscribers", "newPublisher", newPublisher.ID, "newStreamPath", newPublisher.StreamPath)
	var remain SubscriberCollection
	for subscriber := range p.SubscriberRange {
		if subscriber.Type != SubscribeTypeServer {
			remain.Add(subscriber)
		} else {
			newPublisher.AddSubscriber(subscriber)
		}
	}
	p.Subscribers = remain
	p.BufferTime = p.Plugin.GetCommonConf().Publish.BufferTime
	p.AudioTrack.SetMinBuffer(p.BufferTime)
	p.VideoTrack.SetMinBuffer(p.BufferTime)
	if p.State == PublisherStateSubscribed {
		p.State = PublisherStateWaitSubscriber
		if p.DelayCloseTimeout > 0 {
			p.TimeoutTimer.Reset(p.DelayCloseTimeout)
		}
	}
}

func (p *Publisher) takeOver(old *Publisher) {
	if old.HasAudioTrack() {
		p.AudioTrack.baseTs = old.AudioTrack.LastTs
	}
	if old.HasVideoTrack() {
		p.VideoTrack.baseTs = old.VideoTrack.LastTs
	}
	old.Stop(ErrKick)
	p.Info("takeOver", "old", old.ID)
	if old.Subscribers.Length > 0 {
		p.Info(fmt.Sprintf("subscriber +%d", old.Subscribers.Length))
		for subscriber := range old.SubscriberRange {
			subscriber.Publisher = p
			if subscriber.BufferTime > p.BufferTime {
				p.BufferTime = subscriber.BufferTime
			}
		}
	}
	old.AudioTrack.Dispose()
	old.VideoTrack.Dispose()
	old.Subscribers = SubscriberCollection{}
}

func (p *Publisher) WaitTrack() (err error) {
	var v, a = pkg.ErrNoTrack, pkg.ErrNoTrack
	if p.PubVideo {
		v = p.videoReady.Await()
	}
	if p.PubAudio {
		a = p.audioReady.Await()
	}
	if v != nil && a != nil {
		return ErrNoTrack
	}
	return
}

func (p *Publisher) NoVideo() {
	p.PubVideo = false
	if p.videoReady != nil {
		p.videoReady.Reject(ErrMuted)
	}
}

func (p *Publisher) NoAudio() {
	p.PubAudio = false
	if p.audioReady != nil {
		p.audioReady.Reject(ErrMuted)
	}
}

func (p *Publisher) Pause() {
	if p.Paused != nil {
		return
	}
	p.Paused = util.NewPromise(p)
	p.pauseTime = time.Now()
}

func (p *Publisher) Resume() {
	if p.Paused == nil {
		return
	}
	p.Paused.Resolve()
	p.Paused = nil
	p.VideoTrack.pausedTime += time.Since(p.pauseTime)
	p.AudioTrack.pausedTime += time.Since(p.pauseTime)
}

func (p *Publisher) Seek(ts time.Time) {
	p.Info("seek", "time", ts)
	if p.OnSeek != nil {
		p.OnSeek(ts)
	}
}

func (p *Publisher) GetPosition() (t time.Time) {
	if p.OnGetPosition != nil {
		return p.OnGetPosition()
	}
	return
}
