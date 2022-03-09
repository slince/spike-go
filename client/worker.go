package client

import (
	"fmt"
	"github.com/slince/spike/client/proxy"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
	"time"
)

type Worker struct {
	cli *Client
	tun          tunnel.Tunnel
	localAddress string
}

var defaultHost = "127.0.0.1"

func newWorker(cli *Client, tun tunnel.Tunnel) *Worker{
	var localHost = defaultHost
	if len(tun.LocalHost) > 0 {
		localHost = tun.LocalHost
	}
	return &Worker{
		cli:          cli,
		tun:          tun,
		localAddress: net.JoinHostPort(localHost, strconv.Itoa(tun.LocalPort)),
	}
}

func (w *Worker) newLocalConn() (net.Conn, error){
	var con, err = net.DialTimeout("tcp", w.localAddress, 5 * time.Second)
	if err != nil {
		w.cli.logger.Error("Failed to connect local service: ", err)
	} else {
		w.cli.logger.Info("Connected to the local service: ", w.localAddress)
	}
	return con,err
}

func (w *Worker) Start() {
	var proxyConn, err = w.cli.NewConn()
	if err != nil {
		return
	}

	var bridge = transfer.NewBridge(ft, proxyConn, proxyConn)
	_ = bridge.Write(&cmd.RegisterProxy{Tunnel: w.tun, ClientId: w.cli.id})

	handler, err := w.createHandler(proxyConn)
	if err != nil {
		return
	}
	err = handler.Start()
	if err != nil {
		return
	}
	w.cli.logger.Info("The worker is closed")
}

func (w *Worker) createHandler(proxyConn net.Conn) (proxy.Handler, error){
	var handler proxy.Handler
	var err error
	switch w.tun.Protocol {
	case "tcp":
		handler = proxy.NewTcpHandler(w.cli.logger, w.localAddress, proxyConn)
	case "udp":
		handler = proxy.NewUdpHandler(w.cli.logger, w.localAddress, proxyConn)
	default:
		err = fmt.Errorf("unsupported tunel protocol %s", w.tun.Protocol)
	}
	return handler, err
}