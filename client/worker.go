package client

import (
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
)

type Worker struct {
	cli *Client
	tun tunnel.Tunnel
	localConn net.Conn
}


func (w *Worker) start() error {
	var address = "127.0.0.1:" + strconv.Itoa(int(w.tun.LocalPort))
	var conn, err = net.DialTimeout("tcp", address, 5)
	if err != nil {
		w.cli.logger.Warn("Failed to connect local service: ", err)
		return err
	}
	w.cli.logger.Info("Connect to local service successfully ", address)

}