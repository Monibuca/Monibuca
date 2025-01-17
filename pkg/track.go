package pkg

import (
	"context"
	"log/slog"
	"math"
	"reflect"
	"time"

	"m7s.live/v5/pkg/codec"
	"m7s.live/v5/pkg/config"
	"m7s.live/v5/pkg/task"

	"m7s.live/v5/pkg/util"
)

type (
	Track struct {
		*slog.Logger
		ready       *util.Promise
		FrameType   reflect.Type
		bytesIn     int
		frameCount  int
		lastBPSTime time.Time
		BPS         int
		FPS         int
	}

	DataTrack struct {
		Track
	}
	TsTamer struct {
		BaseTs, LastTs time.Duration
	}
	AVTrack struct {
		Track
		*RingWriter
		codec.ICodecCtx
		Allocator     *util.ScalableMemoryAllocator
		SequenceFrame IAVFrame
		WrapIndex     int
		TsTamer
	}
)

func NewAVTrack(args ...any) (t *AVTrack) {
	t = &AVTrack{}
	for _, arg := range args {
		switch v := arg.(type) {
		case IAVFrame:
			t.FrameType = reflect.TypeOf(v)
			t.Allocator = v.GetAllocator()
		case reflect.Type:
			t.FrameType = v
		case *slog.Logger:
			t.Logger = v
		case *AVTrack:
			t.Logger = v.Logger.With("subtrack", t.FrameType.String())
			t.RingWriter = v.RingWriter
			t.ready = util.NewPromiseWithTimeout(context.TODO(), time.Second*5)
		case *config.Publish:
			t.RingWriter = NewRingWriter(v.RingSize)
			t.BufferRange[0] = v.BufferTime
			t.RingWriter.SLogger = t.Logger
		case *util.Promise:
			t.ready = v
		}
	}
	//t.ready = util.NewPromise(struct{}{})
	t.Info("create")
	return
}

func (t *Track) GetKey() reflect.Type {
	return t.FrameType
}

func (t *Track) AddBytesIn(n int) {
	t.bytesIn += n
	t.frameCount++
	if dur := time.Since(t.lastBPSTime); dur > time.Second {
		t.BPS = int(float64(t.bytesIn) / dur.Seconds())
		t.bytesIn = 0
		t.FPS = int(float64(t.frameCount) / dur.Seconds())
		t.frameCount = 0
		t.lastBPSTime = time.Now()
	}
}

func (t *AVTrack) Ready(err error) {
	if t.ready.IsPending() {
		if err != nil {
			t.Error("ready", "err", err)
		} else {
			switch ctx := t.ICodecCtx.(type) {
			case IVideoCodecCtx:
				t.Info("ready", "info", t.ICodecCtx.GetInfo(), "width", ctx.Width(), "height", ctx.Height())
			case IAudioCodecCtx:
				t.Info("ready", "info", t.ICodecCtx.GetInfo(), "channels", ctx.GetChannels(), "sample_rate", ctx.GetSampleRate())
			}
		}
		t.ready.Fulfill(err)
	}
}

func (t *Track) Ready(err error) {
	if t.ready.IsPending() {
		if err != nil {
			t.Error("ready", "err", err)
		} else {
			t.Info("ready")
		}
		t.ready.Fulfill(err)
	}
}

func (t *Track) IsReady() bool {
	return !t.ready.IsPending()
}

func (t *Track) WaitReady() error {
	return t.ready.Await()
}

func (t *Track) Trace(msg string, fields ...any) {
	t.Log(context.TODO(), task.TraceLevel, msg, fields...)
}

func (t *TsTamer) Tame(ts time.Duration, fps int, scale float64) (result time.Duration) {
	if t.LastTs == 0 {
		t.BaseTs -= ts
	}
	result = max(1*time.Millisecond, t.BaseTs+ts)
	if fps > 0 {
		frameDur := float64(time.Second) / float64(fps)
		if math.Abs(float64(result-t.LastTs)) > 10*frameDur { //时间戳突变
			// t.Warn("timestamp mutation", "fps", t.FPS, "lastTs", uint32(t.LastTs/time.Millisecond), "ts", uint32(frame.Timestamp/time.Millisecond), "frameDur", time.Duration(frameDur))
			result = t.LastTs + time.Duration(frameDur)
			t.BaseTs = result - ts
		}
	}
	t.LastTs = result
	result = time.Duration(float64(result) / scale)
	return
}
