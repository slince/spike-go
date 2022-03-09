package client

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
	"time"
)

type Worker struct {
	cli *Client
	tun          tunnel.Tunnel
	LocalAddress string
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
		LocalAddress: net.JoinHostPort(localHost, strconv.Itoa(tun.LocalPort)),
	}
}

func (w *Worker) newLocalConn() (net.Conn, error){
	var con, err = net.DialTimeout("tcp", w.LocalAddress, 5 * time.Second)
	if err != nil {
		w.cli.logger.Error("Failed to connect local service: ", err)
	} else {
		w.cli.logger.Info("Connected to the local service: ", w.LocalAddress)
	}
	return con,err
}

func (w *Worker) start() {
	var proxyConn net.Conn
	var err error

	proxyConn, err = w.cli.NewConn()
	if err != nil {
		return
	}

	var bridge = transfer.NewBridge(ft, proxyConn, proxyConn)
	_ = bridge.Write(&cmd.RegisterProxy{Tunnel: w.tun, ClientId: w.cli.id})

	localConn, err := w.newLocalConn()
	if err != nil {
		return
	}
	conn.Combine(localConn, proxyConn)

	w.cli.logger.Info("The worker is closed")
}