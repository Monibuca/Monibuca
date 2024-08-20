package config

import (
	"context"
	"crypto/tls"
	"log/slog"
	"m7s.live/m7s/v5/pkg/util"

	"github.com/quic-go/quic-go"
)

type QuicConfig interface {
	ListenQuic(context.Context, func(connection quic.Connection)) error
}

type Quic struct {
	ListenAddr string `desc:"监听地址，格式为ip:port，ip 可省略默认监听所有网卡"`
	CertFile   string `desc:"证书文件"`
	KeyFile    string `desc:"私钥文件"`
	AutoListen bool   `default:"true" desc:"是否自动监听"`
}

func (q *Quic) CreateQUICTask(logger *slog.Logger, handler func(connection quic.Connection) util.ITask) *ListenQuicTask {
	ret := &ListenQuicTask{
		Quic:    q,
		handler: handler,
	}
	ret.Logger = logger.With("addr", q.ListenAddr)
	return ret
}

type ListenQuicTask struct {
	util.MarcoLongTask
	*Quic
	*quic.Listener
	handler func(connection quic.Connection) util.ITask
}

func (task *ListenQuicTask) Start() (err error) {
	var ltsc *tls.Config
	ltsc, err = GetTLSConfig(task.CertFile, task.KeyFile)
	if err != nil {
		return
	}
	task.Listener, err = quic.ListenAddr(task.ListenAddr, ltsc, &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		task.Error("listen quic error", err)
		return
	}
	task.Info("listen quic on", task.ListenAddr)
	return
}

func (task *ListenQuicTask) Go() error {
	for {
		conn, err := task.Accept(task.Context)
		if err != nil {
			return err
		}
		subTask := task.handler(conn)
		task.AddTask(subTask)
	}
}

func (task *ListenQuicTask) Dispose() {
	_ = task.Listener.Close()
}
