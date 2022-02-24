package client

import (
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
)

type Worker struct {
	cli *Client
	tun tunnel.Tunnel
}

func newWorker(cli *Client, tun tunnel.Tunnel) *Worker{
	return &Worker{
		cli: cli,
		tun: tun,
	}
}

func (w *Worker) start() {
	var address = "127.0.0.1:" + strconv.Itoa(int(w.tun.LocalPort))
	var localConn, err = net.DialTimeout("tcp", address, 5)
	if err != nil {
		w.cli.logger.Warn("Failed to connect local service: ", err)
		return
	}
	w.cli.logger.Info("Connect to local service successfully ", address)
	var proxyConn net.Conn
	proxyConn, err = w.cli.newConn()
	conn.Combine(localConn, proxyConn, func(alive net.Conn) {
		_ = localConn.Close()
		_ = proxyConn.Close()
	})
}