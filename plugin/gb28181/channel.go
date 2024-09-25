package plugin_gb28181

import (
	"log/slog"
	"sync/atomic"
	"time"

	"m7s.live/m7s/v5"
	"m7s.live/m7s/v5/pkg/config"
	"m7s.live/m7s/v5/pkg/util"
	gb28181 "m7s.live/m7s/v5/plugin/gb28181/pkg"
)

type RecordRequest struct {
	SN, SumNum int
	Response   []gb28181.Record
	*util.Promise
}

func (r *RecordRequest) GetKey() int {
	return r.SN
}

type Channel struct {
	Device              *Device      // 所属设备
	State               atomic.Int32 // 通道状态,0:空闲,1:正在invite,2:正在播放/对讲
	GpsTime             time.Time    // gps时间
	Longitude, Latitude string       // 经度
	RecordReqs          util.Collection[int, *RecordRequest]
	*slog.Logger
	gb28181.ChannelInfo
	AbstractDevice *m7s.Device
}

func (c *Channel) GetKey() string {
	return c.DeviceID
}

func (c *Channel) Pull() {
	c.Device.plugin.Pull(c.AbstractDevice.GetStreamPath(), config.Pull{URL: c.AbstractDevice.PullURL,MaxRetry: -1})
}
