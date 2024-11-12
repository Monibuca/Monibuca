package plugin_transcode

import (
	m7s "m7s.live/pro"
	"m7s.live/pro/plugin/transcode/pb"
	transcode "m7s.live/pro/plugin/transcode/pkg"
)

var _ = m7s.InstallPlugin[TranscodePlugin](transcode.NewTransform, pb.RegisterApiHandler, &pb.Api_ServiceDesc)

type TranscodePlugin struct {
	pb.UnimplementedApiServer
	m7s.Plugin
	LogToFile string
}
