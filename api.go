package m7s

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"m7s.live/v5/pkg/task"

	myip "github.com/husanpao/ip"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	. "github.com/shirou/gopsutil/v4/net"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"
	"m7s.live/v5/pb"
	"m7s.live/v5/pkg"
	"m7s.live/v5/pkg/config"
	"m7s.live/v5/pkg/util"
)

var localIP string

func (s *Server) SysInfo(context.Context, *emptypb.Empty) (res *pb.SysInfoResponse, err error) {
	if localIP == "" {
		localIP = myip.LocalIP()
		// if conn, err := net.Dial("udp", "114.114.114.114:80"); err == nil {
		// 	localIP, _, _ = strings.Cut(conn.LocalAddr().String(), ":")
		// }
	}
	res = &pb.SysInfoResponse{
		Code:    0,
		Message: "success",
		Data: &pb.SysInfoData{
			Version:   Version,
			LocalIP:   localIP,
			PublicIP:  util.GetPublicIP(""),
			StartTime: timestamppb.New(s.StartTime),
			GoVersion: runtime.Version(),
			Os:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Cpus:      int32(runtime.NumCPU()),
		},
	}
	for p := range s.Plugins.Range {
		res.Data.Plugins = append(res.Data.Plugins, &pb.PluginInfo{
			Name:        p.Meta.Name,
			PushAddr:    p.PushAddr,
			PlayAddr:    p.PlayAddr,
			Description: p.GetDescriptions(),
		})
	}
	return
}

func (s *Server) DisabledPlugins(ctx context.Context, _ *emptypb.Empty) (res *pb.DisabledPluginsResponse, err error) {
	res = &pb.DisabledPluginsResponse{
		Data: make([]*pb.PluginInfo, len(s.disabledPlugins)),
	}
	for i, p := range s.disabledPlugins {
		res.Data[i] = &pb.PluginInfo{
			Name:        p.Meta.Name,
			Description: p.GetDescriptions(),
		}
	}
	return
}

// /api/stream/annexb/{streamPath}
func (s *Server) api_Stream_AnnexB_(rw http.ResponseWriter, r *http.Request) {
	publisher, ok := s.Streams.Get(r.PathValue("streamPath"))
	if !ok || publisher.VideoTrack.AVTrack == nil {
		http.Error(rw, pkg.ErrNotFound.Error(), http.StatusNotFound)
		return
	}
	err := publisher.VideoTrack.WaitReady()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/octet-stream")
	reader := pkg.NewAVRingReader(publisher.VideoTrack.AVTrack, "Origin")
	err = reader.StartRead(publisher.VideoTrack.GetIDR())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.StopRead()
	if reader.Value.Raw == nil {
		if err = reader.Value.Demux(publisher.VideoTrack.ICodecCtx); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var annexb pkg.AnnexB
	var t pkg.AVTrack

	t.ICodecCtx, t.SequenceFrame, err = annexb.ConvertCtx(publisher.VideoTrack.ICodecCtx)
	if t.ICodecCtx == nil {
		http.Error(rw, "unsupported codec", http.StatusInternalServerError)
		return
	}
	annexb.Mux(t.ICodecCtx, &reader.Value)
	_, err = annexb.WriteTo(rw)
}

func (s *Server) getStreamInfo(pub *Publisher) (res *pb.StreamInfoResponse, err error) {
	tmp, _ := json.Marshal(pub.GetDescriptions())
	res = &pb.StreamInfoResponse{
		Data: &pb.StreamInfo{
			Meta:      string(tmp),
			Path:      pub.StreamPath,
			State:     int32(pub.State),
			StartTime: timestamppb.New(pub.StartTime),
			// Subscribers: int32(pub.Subscribers.Length),
			PluginName: pub.Plugin.Meta.Name,
			Type:       pub.Type,
			Speed:      float32(pub.Speed),
			StopOnIdle: pub.DelayCloseTimeout > 0,
			IsPaused:   pub.Paused != nil,
			Gop:        int32(pub.GOP),
			BufferTime: durationpb.New(pub.BufferTime),
		},
	}
	var audioBpsOut, videoBpsOut uint32
	var serverSubCount int32
	for sub := range pub.Subscribers.Range {
		if sub.AudioReader != nil {
			audioBpsOut += sub.AudioReader.BPS
		}
		if sub.VideoReader != nil {
			videoBpsOut += sub.VideoReader.BPS
		}
		if sub.Type == SubscribeTypeServer {
			serverSubCount++
		}
	}
	res.Data.Subscribers = serverSubCount
	if t := pub.AudioTrack.AVTrack; t != nil {
		if t.ICodecCtx != nil {
			res.Data.AudioTrack = &pb.AudioTrackInfo{
				Codec:  t.FourCC().String(),
				Meta:   t.GetInfo(),
				Bps:    uint32(t.BPS),
				BpsOut: audioBpsOut,
				Fps:    uint32(t.FPS),
				Delta:  pub.AudioTrack.Delta.String(),
			}
			res.Data.AudioTrack.SampleRate = uint32(t.ICodecCtx.(pkg.IAudioCodecCtx).GetSampleRate())
			res.Data.AudioTrack.Channels = uint32(t.ICodecCtx.(pkg.IAudioCodecCtx).GetChannels())
		}
	}
	if t := pub.VideoTrack.AVTrack; t != nil {
		if t.ICodecCtx != nil {
			res.Data.VideoTrack = &pb.VideoTrackInfo{
				Codec:  t.FourCC().String(),
				Meta:   t.GetInfo(),
				Bps:    uint32(t.BPS),
				BpsOut: videoBpsOut,
				Fps:    uint32(t.FPS),
				Delta:  pub.VideoTrack.Delta.String(),
				Gop:    uint32(pub.GOP),
			}
			res.Data.VideoTrack.Width = uint32(t.ICodecCtx.(pkg.IVideoCodecCtx).Width())
			res.Data.VideoTrack.Height = uint32(t.ICodecCtx.(pkg.IVideoCodecCtx).Height())
		}
	}
	return
}

func (s *Server) StreamInfo(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.StreamInfoResponse, err error) {
	var recordings []*pb.RecordingDetail
	s.Records.Call(func() error {
		for record := range s.Records.Range {
			if record.StreamPath == req.StreamPath {
				recordings = append(recordings, &pb.RecordingDetail{
					FilePath:   record.FilePath,
					Mode:       record.Mode,
					Fragment:   durationpb.New(record.Fragment),
					Append:     record.Append,
					PluginName: record.Plugin.Meta.Name,
				})
			}
		}
		return nil
	})
	s.Streams.Call(func() error {
		if pub, ok := s.Streams.Get(req.StreamPath); ok {
			res, err = s.getStreamInfo(pub)
			if err != nil {
				return err
			}
			res.Data.Recording = recordings
		} else {
			err = pkg.ErrNotFound
		}
		return nil
	})
	return
}

func (s *Server) TaskTree(context.Context, *emptypb.Empty) (res *pb.TaskTreeResponse, err error) {
	var fillData func(m task.ITask) *pb.TaskTreeData
	fillData = func(m task.ITask) (res *pb.TaskTreeData) {
		if m == nil {
			return
		}
		t := m.GetTask()
		res = &pb.TaskTreeData{
			Id:          m.GetTaskID(),
			Pointer:     uint64(t.GetTaskPointer()),
			State:       uint32(m.GetState()),
			Type:        uint32(m.GetTaskType()),
			Owner:       m.GetOwnerType(),
			StartTime:   timestamppb.New(t.StartTime),
			Description: m.GetDescriptions(),
			StartReason: t.StartReason,
		}
		if job, ok := m.(task.IJob); ok {
			if blockedTask := job.Blocked(); blockedTask != nil {
				res.Blocked = fillData(blockedTask)
			}
			for t := range job.RangeSubTask {
				child := fillData(t)
				if child == nil {
					continue
				}
				res.Children = append(res.Children, child)
			}
		}
		return
	}
	res = &pb.TaskTreeResponse{Data: fillData(&Servers)}
	return
}

func (s *Server) StopTask(ctx context.Context, req *pb.RequestWithId64) (resp *pb.SuccessResponse, err error) {
	t := task.FromPointer(uintptr(req.Id))
	if t == nil {
		return nil, pkg.ErrNotFound
	}
	t.Stop(task.ErrStopByUser)
	return &pb.SuccessResponse{}, nil
}

func (s *Server) RestartTask(ctx context.Context, req *pb.RequestWithId64) (resp *pb.SuccessResponse, err error) {
	t := task.FromPointer(uintptr(req.Id))
	if t == nil {
		return nil, pkg.ErrNotFound
	}
	t.Stop(task.ErrRestart)
	return &pb.SuccessResponse{}, nil
}

func (s *Server) GetRecording(ctx context.Context, req *emptypb.Empty) (resp *pb.RecordingListResponse, err error) {
	s.Records.Call(func() error {
		resp = &pb.RecordingListResponse{}
		for record := range s.Records.Range {
			resp.Data = append(resp.Data, &pb.Recording{
				StreamPath: record.StreamPath,
				StartTime:  timestamppb.New(record.StartTime),
				Type:       reflect.TypeOf(record.recorder).String(),
				Pointer:    uint64(record.GetTaskPointer()),
			})
		}
		return nil
	})
	return
}

func (s *Server) GetSubscribers(context.Context, *pb.SubscribersRequest) (res *pb.SubscribersResponse, err error) {
	s.Streams.Call(func() error {
		var subscribers []*pb.SubscriberSnapShot
		for subscriber := range s.Subscribers.Range {
			meta, _ := json.Marshal(subscriber.GetDescriptions())
			snap := &pb.SubscriberSnapShot{
				Id:         subscriber.ID,
				StartTime:  timestamppb.New(subscriber.StartTime),
				Meta:       string(meta),
				Type:       subscriber.Type,
				PluginName: subscriber.Plugin.Meta.Name,
				SubMode:    int32(subscriber.SubMode),
				SyncMode:   int32(subscriber.SyncMode),
				BufferTime: durationpb.New(subscriber.BufferTime),
				RemoteAddr: subscriber.RemoteAddr,
			}
			if ar := subscriber.AudioReader; ar != nil {
				snap.AudioReader = &pb.RingReaderSnapShot{
					Sequence:  ar.Value.Sequence,
					Timestamp: ar.AbsTime,
					Delay:     ar.Delay,
					State:     int32(ar.State),
					Bps:       ar.BPS,
				}
			}
			if vr := subscriber.VideoReader; vr != nil {
				snap.VideoReader = &pb.RingReaderSnapShot{
					Sequence:  vr.Value.Sequence,
					Timestamp: vr.AbsTime,
					Delay:     vr.Delay,
					State:     int32(vr.State),
					Bps:       vr.BPS,
				}
			}
			subscribers = append(subscribers, snap)
		}
		res = &pb.SubscribersResponse{
			Data:  subscribers,
			Total: int32(s.Subscribers.Length),
		}
		return nil
	})
	return
}
func (s *Server) AudioTrackSnap(_ context.Context, req *pb.StreamSnapRequest) (res *pb.TrackSnapShotResponse, err error) {
	s.Streams.Call(func() error {
		if pub, ok := s.Streams.Get(req.StreamPath); ok && pub.HasAudioTrack() {
			data := &pb.TrackSnapShotData{}
			if pub.AudioTrack.Allocator != nil {
				for _, memlist := range pub.AudioTrack.Allocator.GetChildren() {
					var list []*pb.MemoryBlock
					for _, block := range memlist.GetBlocks() {
						list = append(list, &pb.MemoryBlock{
							S: uint32(block.Start),
							E: uint32(block.End),
						})
					}
					data.Memory = append(data.Memory, &pb.MemoryBlockGroup{List: list, Size: uint32(memlist.Size)})
				}
			}
			pub.AudioTrack.Ring.Do(func(v *pkg.AVFrame) {
				if len(v.Wraps) > 0 {
					var snap pb.TrackSnapShot
					snap.Sequence = v.Sequence
					snap.Timestamp = uint32(v.Timestamp / time.Millisecond)
					snap.WriteTime = timestamppb.New(v.WriteTime)
					snap.Wrap = make([]*pb.Wrap, len(v.Wraps))
					snap.KeyFrame = v.IDR
					data.RingDataSize += uint32(v.Wraps[0].GetSize())
					for i, wrap := range v.Wraps {
						snap.Wrap[i] = &pb.Wrap{
							Timestamp: uint32(wrap.GetTimestamp() / time.Millisecond),
							Size:      uint32(wrap.GetSize()),
							Data:      wrap.String(),
						}
					}
					data.Ring = append(data.Ring, &snap)
				}
			})
			res = &pb.TrackSnapShotResponse{
				Code:    0,
				Message: "success",
				Data:    data,
			}
		} else {
			err = pkg.ErrNotFound
		}
		return nil
	})
	return
}
func (s *Server) api_VideoTrack_SSE(rw http.ResponseWriter, r *http.Request) {
	streamPath := r.PathValue("streamPath")
	if r.URL.RawQuery != "" {
		streamPath += "?" + r.URL.RawQuery
	}
	suber, err := s.SubscribeWithConfig(r.Context(), streamPath, config.Subscribe{
		SubVideo: true,
		SubType:  SubscribeTypeAPI,
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	sse := util.NewSSE(rw, r.Context())
	PlayBlock(suber, (func(frame *pkg.AVFrame) (err error))(nil), func(frame *pkg.AVFrame) (err error) {
		var snap pb.TrackSnapShot
		snap.Sequence = frame.Sequence
		snap.Timestamp = uint32(frame.Timestamp / time.Millisecond)
		snap.WriteTime = timestamppb.New(frame.WriteTime)
		snap.Wrap = make([]*pb.Wrap, len(frame.Wraps))
		snap.KeyFrame = frame.IDR
		for i, wrap := range frame.Wraps {
			snap.Wrap[i] = &pb.Wrap{
				Timestamp: uint32(wrap.GetTimestamp() / time.Millisecond),
				Size:      uint32(wrap.GetSize()),
				Data:      wrap.String(),
			}
		}
		return sse.WriteJSON(&snap)
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) api_AudioTrack_SSE(rw http.ResponseWriter, r *http.Request) {
	streamPath := r.PathValue("streamPath")
	if r.URL.RawQuery != "" {
		streamPath += "?" + r.URL.RawQuery
	}
	suber, err := s.SubscribeWithConfig(r.Context(), streamPath, config.Subscribe{
		SubAudio: true,
		SubType:  SubscribeTypeAPI,
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	sse := util.NewSSE(rw, r.Context())
	PlayBlock(suber, func(frame *pkg.AVFrame) (err error) {
		var snap pb.TrackSnapShot
		snap.Sequence = frame.Sequence
		snap.Timestamp = uint32(frame.Timestamp / time.Millisecond)
		snap.WriteTime = timestamppb.New(frame.WriteTime)
		snap.Wrap = make([]*pb.Wrap, len(frame.Wraps))
		snap.KeyFrame = frame.IDR
		for i, wrap := range frame.Wraps {
			snap.Wrap[i] = &pb.Wrap{
				Timestamp: uint32(wrap.GetTimestamp() / time.Millisecond),
				Size:      uint32(wrap.GetSize()),
				Data:      wrap.String(),
			}
		}
		return sse.WriteJSON(&snap)
	}, (func(frame *pkg.AVFrame) (err error))(nil))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) VideoTrackSnap(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.TrackSnapShotResponse, err error) {
	s.Streams.Call(func() error {
		if pub, ok := s.Streams.Get(req.StreamPath); ok && pub.HasVideoTrack() {
			data := &pb.TrackSnapShotData{}
			if pub.VideoTrack.Allocator != nil {
				for _, memlist := range pub.VideoTrack.Allocator.GetChildren() {
					var list []*pb.MemoryBlock
					for _, block := range memlist.GetBlocks() {
						list = append(list, &pb.MemoryBlock{
							S: uint32(block.Start),
							E: uint32(block.End),
						})
					}
					data.Memory = append(data.Memory, &pb.MemoryBlockGroup{List: list, Size: uint32(memlist.Size)})
				}
			}
			pub.VideoTrack.Ring.Do(func(v *pkg.AVFrame) {
				if len(v.Wraps) > 0 {
					var snap pb.TrackSnapShot
					snap.Sequence = v.Sequence
					snap.Timestamp = uint32(v.Timestamp / time.Millisecond)
					snap.WriteTime = timestamppb.New(v.WriteTime)
					snap.Wrap = make([]*pb.Wrap, len(v.Wraps))
					snap.KeyFrame = v.IDR
					data.RingDataSize += uint32(v.Wraps[0].GetSize())
					for i, wrap := range v.Wraps {
						snap.Wrap[i] = &pb.Wrap{
							Timestamp: uint32(wrap.GetTimestamp() / time.Millisecond),
							Size:      uint32(wrap.GetSize()),
							Data:      wrap.String(),
						}
					}
					data.Ring = append(data.Ring, &snap)
				}
			})
			res = &pb.TrackSnapShotResponse{
				Code:    0,
				Message: "success",
				Data:    data,
			}
		} else {
			err = pkg.ErrNotFound
		}
		return nil
	})
	return
}

// Restart stops the server with a restart error and returns
// a success response. This method is used to restart the server
// gracefully.
func (s *Server) Restart(ctx context.Context, req *pb.RequestWithId) (res *pb.SuccessResponse, err error) {
	s.Stop(pkg.ErrRestart)
	return &pb.SuccessResponse{}, err
}

func (s *Server) Shutdown(ctx context.Context, req *pb.RequestWithId) (res *pb.SuccessResponse, err error) {
	s.Stop(task.ErrStopByUser)
	return &pb.SuccessResponse{}, err
}

func (s *Server) ChangeSubscribe(ctx context.Context, req *pb.ChangeSubscribeRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if subscriber, ok := s.Subscribers.Get(req.Id); ok {
			if pub, ok := s.Streams.Get(req.StreamPath); ok {
				subscriber.Publisher.RemoveSubscriber(subscriber)
				subscriber.StreamPath = req.StreamPath
				pub.AddSubscriber(subscriber)
				return nil
			}
		}
		err = pkg.ErrNotFound
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) StopSubscribe(ctx context.Context, req *pb.RequestWithId) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if subscriber, ok := s.Subscribers.Get(req.Id); ok {
			subscriber.Stop(errors.New("stop by api"))
		} else {
			err = pkg.ErrNotFound
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) PauseStream(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if s, ok := s.Streams.Get(req.StreamPath); ok {
			s.Pause()
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) ResumeStream(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if s, ok := s.Streams.Get(req.StreamPath); ok {
			s.Resume()
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) SetStreamSpeed(ctx context.Context, req *pb.SetStreamSpeedRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if s, ok := s.Streams.Get(req.StreamPath); ok {
			s.Speed = float64(req.Speed)
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) SeekStream(ctx context.Context, req *pb.SeekStreamRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if s, ok := s.Streams.Get(req.StreamPath); ok {
			s.Seek(time.Unix(int64(req.TimeStamp), 0))
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

func (s *Server) StopPublish(ctx context.Context, req *pb.StreamSnapRequest) (res *pb.SuccessResponse, err error) {
	s.Streams.Call(func() error {
		if s, ok := s.Streams.Get(req.StreamPath); ok {
			s.Stop(task.ErrStopByUser)
		}
		return nil
	})
	return &pb.SuccessResponse{}, err
}

// /api/stream/list
func (s *Server) StreamList(_ context.Context, req *pb.StreamListRequest) (res *pb.StreamListResponse, err error) {
	recordingMap := make(map[string][]*pb.RecordingDetail)
	s.Records.Call(func() error {
		for record := range s.Records.Range {
			recordingMap[record.StreamPath] = append(recordingMap[record.StreamPath], &pb.RecordingDetail{
				FilePath:   record.FilePath,
				Mode:       record.Mode,
				Fragment:   durationpb.New(record.Fragment),
				Append:     record.Append,
				PluginName: record.Plugin.Meta.Name,
				Pointer:    uint64(record.GetTaskPointer()),
			})
		}
		return nil
	})
	s.Streams.Call(func() error {
		var streams []*pb.StreamInfo
		for publisher := range s.Streams.Range {
			info, err := s.getStreamInfo(publisher)
			if err != nil {
				continue
			}
			info.Data.Recording = recordingMap[info.Data.Path]
			streams = append(streams, info.Data)
		}
		res = &pb.StreamListResponse{Data: streams, Total: int32(s.Streams.Length), PageNum: req.PageNum, PageSize: req.PageSize}
		return nil
	})
	return
}

func (s *Server) WaitList(context.Context, *emptypb.Empty) (res *pb.StreamWaitListResponse, err error) {
	s.Streams.Call(func() error {
		res = &pb.StreamWaitListResponse{
			List: make(map[string]int32),
		}
		for subs := range s.Waiting.Range {
			res.List[subs.StreamPath] = int32(subs.Length)
		}
		return nil
	})
	return
}

func (s *Server) Api_Summary_SSE(rw http.ResponseWriter, r *http.Request) {
	util.ReturnFetchValue(func() *pb.SummaryResponse {
		ret, _ := s.Summary(r.Context(), nil)
		return ret
	}, rw, r)
}

func (s *Server) Api_Stream_Position_SSE(rw http.ResponseWriter, r *http.Request) {
	streamPath := r.URL.Query().Get("streamPath")
	util.ReturnFetchValue(func() (t time.Time) {
		s.Streams.Call(func() error {
			if pub, ok := s.Streams.Get(streamPath); ok {
				t = pub.GetPosition()
			}
			return nil
		})
		return
	}, rw, r)
}

// func (s *Server) Api_Vod_Position(rw http.ResponseWriter, r *http.Request) {
// 	streamPath := r.URL.Query().Get("streamPath")
// 	s.Streams.Call(func() error {
// 		if pub, ok := s.Streams.Get(streamPath); ok {
// 			t = pub.GetPosition()
// 		}
// 		return nil
// 	})
// }

func (s *Server) Summary(context.Context, *emptypb.Empty) (res *pb.SummaryResponse, err error) {
	dur := time.Since(s.lastSummaryTime)
	if dur < time.Second {
		res = s.lastSummary
		return
	}
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	nv, _ := IOCounters(true)
	res = &pb.SummaryResponse{
		Memory: &pb.Usage{
			Total: v.Total >> 20,
			Free:  v.Available >> 20,
			Used:  v.Used >> 20,
			Usage: float32(v.UsedPercent),
		},
		HardDisk: &pb.Usage{
			Total: d.Total >> 30,
			Free:  d.Free >> 30,
			Used:  d.Used >> 30,
			Usage: float32(d.UsedPercent),
		},
	}
	if cc, _ := cpu.Percent(0, false); len(cc) > 0 {
		res.CpuUsage = float32(cc[0])
	}
	netWorks := []*pb.NetWorkInfo{}
	for i, n := range nv {
		info := &pb.NetWorkInfo{
			Name:    n.Name,
			Receive: n.BytesRecv,
			Sent:    n.BytesSent,
		}
		if s.lastSummary != nil && len(s.lastSummary.NetWork) > i {
			info.ReceiveSpeed = (n.BytesRecv - s.lastSummary.NetWork[i].Receive) / uint64(dur.Seconds())
			info.SentSpeed = (n.BytesSent - s.lastSummary.NetWork[i].Sent) / uint64(dur.Seconds())
		}
		netWorks = append(netWorks, info)
	}
	res.StreamCount = int32(s.Streams.Length)
	res.PullCount = int32(s.Pulls.Length)
	res.PushCount = int32(s.Pushs.Length)
	res.SubscribeCount = int32(s.Subscribers.Length)
	res.RecordCount = int32(s.Records.Length)
	res.TransformCount = int32(s.Transforms.Length)
	res.NetWork = netWorks
	s.lastSummary = res
	s.lastSummaryTime = time.Now()
	return
}

// /api/config/json/{name}
func (s *Server) api_Config_JSON_(rw http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var conf *config.Config
	if name == "global" {
		conf = &s.Config
	} else {
		p, ok := s.Plugins.Get(name)
		if !ok {
			http.Error(rw, pkg.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		conf = &p.Config
	}
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(conf.GetMap())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetConfigFile(_ context.Context, req *emptypb.Empty) (res *pb.GetConfigFileResponse, err error) {
	res = &pb.GetConfigFileResponse{}
	res.Data = string(s.configFileContent)
	return
}

func (s *Server) UpdateConfigFile(_ context.Context, req *pb.UpdateConfigFileRequest) (res *pb.SuccessResponse, err error) {
	if s.configFileContent != nil {
		s.configFileContent = []byte(req.Content)
		os.WriteFile(filepath.Join(ExecDir, s.conf.(string)), s.configFileContent, 0644)
		res = &pb.SuccessResponse{}
	} else {
		err = pkg.ErrNotFound
	}
	return
}

func (s *Server) GetConfig(_ context.Context, req *pb.GetConfigRequest) (res *pb.GetConfigResponse, err error) {
	res = &pb.GetConfigResponse{
		Data: &pb.ConfigData{},
	}
	var conf *config.Config
	if req.Name == "global" {
		conf = &s.Config
	} else {
		p, ok := s.Plugins.Get(req.Name)
		if !ok {
			err = pkg.ErrNotFound
			return
		}
		conf = &p.Config
	}
	var mm []byte
	mm, err = yaml.Marshal(conf.File)
	if err != nil {
		return
	}
	res.Data.File = string(mm)

	mm, err = yaml.Marshal(conf.Modify)
	if err != nil {
		return
	}
	res.Data.Modified = string(mm)

	mm, err = yaml.Marshal(conf.GetMap())
	if err != nil {
		return
	}
	res.Data.Merged = string(mm)
	return
}

func (s *Server) ModifyConfig(_ context.Context, req *pb.ModifyConfigRequest) (res *pb.SuccessResponse, err error) {
	var conf *config.Config
	if req.Name == "global" {
		conf = &s.Config
		defer s.SaveConfig()
	} else {
		p, ok := s.Plugins.Get(req.Name)
		if !ok {
			err = pkg.ErrNotFound
			return
		}
		defer p.SaveConfig()
		conf = &p.Config
	}
	var modified map[string]any
	err = yaml.Unmarshal([]byte(req.Yaml), &modified)
	if err != nil {
		return
	}
	conf.ParseModifyFile(modified)
	return
}

func (s *Server) GetRecordList(ctx context.Context, req *pb.ReqRecordList) (resp *pb.ResponseList, err error) {
	if s.DB == nil {
		err = pkg.ErrNoDB
		return
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	offset := (req.PageNum - 1) * req.PageSize // 计算偏移量
	var totalCount int64                       //总条数

	var result []*RecordStream
	query := s.DB.Model(&RecordStream{})
	if strings.Contains(req.StreamPath, "*") {
		query = query.Where("stream_path like ?", strings.ReplaceAll(req.StreamPath, "*", "%"))
	} else if req.StreamPath != "" {
		query = query.Where("stream_path = ?", req.StreamPath)
	}
	if req.Mode != "" {
		query = query.Where("mode = ?", req.Mode)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	startTime, endTime, err := util.TimeRangeQueryParse(url.Values{"range": []string{req.Range}, "start": []string{req.Start}, "end": []string{req.End}})
	if err != nil {
		return
	}
	if !startTime.IsZero() {
		query = query.Where("start_time >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("end_time <= ?", endTime)
	}
	query.Count(&totalCount)
	err = query.Offset(int(offset)).Limit(int(req.PageSize)).Order("start_time desc").Find(&result).Error
	if err != nil {
		return
	}
	resp = &pb.ResponseList{
		TotalCount: uint32(totalCount),
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
	}
	for _, recordFile := range result {
		resp.Data = append(resp.Data, &pb.RecordFile{
			Id:         uint32(recordFile.ID),
			StartTime:  timestamppb.New(recordFile.StartTime),
			EndTime:    timestamppb.New(recordFile.EndTime),
			FilePath:   recordFile.FilePath,
			StreamPath: recordFile.StreamPath,
		})
	}
	return
}

func (s *Server) GetRecordCatalog(ctx context.Context, req *pb.ReqRecordCatalog) (resp *pb.ResponseCatalog, err error) {
	if s.DB == nil {
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
	query := s.DB.Model(&RecordStream{})
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	err = query.Select("stream_path,count(id) as count,min(start_time) as start_time,max(end_time) as end_time").Group("stream_path").Find(&result).Error
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

func (s *Server) DeleteRecord(ctx context.Context, req *pb.ReqRecordDelete) (resp *pb.ResponseDelete, err error) {
	if s.DB == nil {
		err = pkg.ErrNoDB
		return
	}
	ids := req.GetIds()
	var result []*RecordStream
	if len(ids) > 0 {
		s.DB.Find(&result, "stream_path=? AND type=? AND id IN ?", req.StreamPath, req.Type, ids)
	} else {
		startTime, endTime, err := util.TimeRangeQueryParse(url.Values{"range": []string{req.Range}, "start": []string{req.StartTime}, "end": []string{req.EndTime}})
		if err != nil {
			return nil, err
		}
		s.DB.Find(&result, "stream_path=? AND type=? AND start_time>=? AND end_time<=?", req.StreamPath, req.Type, startTime, endTime)
	}
	err = s.DB.Delete(result).Error
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

func (s *Server) GetTransformList(ctx context.Context, req *emptypb.Empty) (res *pb.TransformListResponse, err error) {
	res = &pb.TransformListResponse{}
	s.Transforms.Call(func() error {
		for transform := range s.Transforms.Range {
			info := &pb.Transform{
				StreamPath: transform.StreamPath,
				Target:     transform.Target,
			}
			if transform.TransformJob != nil {
				info.PluginName = transform.TransformJob.Plugin.Meta.Name
				var result []byte
				result, err = yaml.Marshal(transform.TransformJob.Config)
				if err != nil {
					s.Error("marshal transform config failed", "error", err)
					return err
				}
				info.Config = string(result)
			}
			res.Data = append(res.Data, info)
		}
		return nil
	})
	return
}
