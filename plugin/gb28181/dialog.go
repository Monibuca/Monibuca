package plugin_gb28181

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	m7s "m7s.live/v5"
	"m7s.live/v5/pkg/task"
	"m7s.live/v5/pkg/util"
	gb28181 "m7s.live/v5/plugin/gb28181/pkg"
)

type Dialog struct {
	task.Job
	Channel *Channel
	gb28181.InviteOptions
	gb      *GB28181Plugin
	session *sipgo.DialogClientSession
	pullCtx m7s.PullJob
}

func (d *Dialog) GetCallID() string {
	return d.session.InviteRequest.CallID().Value()
}

func (d *Dialog) GetPullJob() *m7s.PullJob {
	return &d.pullCtx
}

func (d *Dialog) Start() (err error) {
	if !d.IsLive() {
		d.pullCtx.PublishConfig.PubType = m7s.PublishTypeVod
	}
	err = d.pullCtx.Publish()
	if err != nil {
		return
	}
	sss := strings.Split(d.pullCtx.RemoteURL, "/")
	deviceId, channelId := sss[0], sss[1]
	if len(sss) == 2 {
		if device, ok := d.gb.devices.Get(deviceId); ok {
			if channel, ok := device.channels.Get(channelId); ok {
				d.Channel = channel
			} else {
				return fmt.Errorf("channel %s not found", channelId)
			}
		} else {
			return fmt.Errorf("device %s not found", deviceId)
		}
	} else if len(sss) == 3 {
		var recordRange util.Range[int]
		err = recordRange.Resolve(sss[2])
	}
	ssrc := d.CreateSSRC(d.gb.Serial)
	d.gb.dialogs.Set(d)
	defer d.gb.dialogs.Remove(d)
	if d.gb.MediaPort.Valid() {
		select {
		case d.MediaPort = <-d.gb.tcpPorts:
			defer func() {
				d.gb.tcpPorts <- d.MediaPort
			}()
		default:
			return fmt.Errorf("no available tcp port")
		}
	} else {
		d.MediaPort = d.gb.MediaPort[0]
	}
	sdpInfo := []string{
		"v=0",
		fmt.Sprintf("o=%s 0 0 IN IP4 %s", d.Channel.DeviceID, d.Channel.Device.mediaIp),
		"s=" + util.Conditional(d.IsLive(), "Play", "Playback"),
		"u=" + d.Channel.DeviceID + ":0",
		"c=IN IP4 " + d.Channel.Device.mediaIp,
		d.String(),
		fmt.Sprintf("m=video %d TCP/RTP/AVP 96", d.MediaPort),
		"a=recvonly",
		"a=rtpmap:96 PS/90000",
		"a=setup:passive",
		"a=connection:new",
		"y=" + ssrc,
	}
	contentTypeHeader := sip.ContentTypeHeader("application/sdp")
	fromHeader := d.Channel.Device.fromHDR
	subjectHeader := sip.NewHeader("Subject", fmt.Sprintf("%s:%s,%s:0", d.Channel.DeviceID, ssrc, d.gb.Serial))
	d.session, err = d.Channel.Device.dialogClient.Invite(d.gb, d.Channel.Device.Recipient, []byte(strings.Join(sdpInfo, "\r\n")+"\r\n"), &contentTypeHeader, subjectHeader, &fromHeader, sip.NewHeader("Allow", "INVITE, ACK, CANCEL, REGISTER, MESSAGE, NOTIFY, BYE"))
	return
}

func (d *Dialog) Run() (err error) {
	err = d.session.WaitAnswer(d.gb, sipgo.AnswerOptions{})
	if err != nil {
		return
	}
	inviteResponseBody := string(d.session.InviteResponse.Body())
	d.Channel.Info("inviteResponse", "body", inviteResponseBody)
	ds := strings.Split(inviteResponseBody, "\r\n")
	for _, l := range ds {
		if ls := strings.Split(l, "="); len(ls) > 1 {
			if ls[0] == "y" && len(ls[1]) > 0 {
				if _ssrc, err := strconv.ParseInt(ls[1], 10, 0); err == nil {
					d.SSRC = uint32(_ssrc)
				} else {
					d.gb.Error("read invite response y ", "err", err)
				}
				//	break
			}
			if ls[0] == "m" && len(ls[1]) > 0 {
				netinfo := strings.Split(ls[1], " ")
				if strings.ToUpper(netinfo[2]) == "TCP/RTP/AVP" {
					d.gb.Debug("device support tcp")
				} else {
					return fmt.Errorf("device not support tcp")
				}
			}
		}
	}
	err = d.session.Ack(d.gb)
	pub := gb28181.NewPSPublisher(d.pullCtx.Publisher)
	pub.Receiver.ListenAddr = fmt.Sprintf(":%d", d.MediaPort)
	d.AddTask(&pub.Receiver)
	pub.Demux()
	return
}

func (d *Dialog) GetKey() uint32 {
	return d.SSRC
}

func (d *Dialog) Dispose() {
	d.session.Close()
}
