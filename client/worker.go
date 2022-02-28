package client

import (
	"fmt"
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
	tun tunnel.Tunnel
	localAddress string
}

var defaultHost = "127.0.0.1"

func newWorker(cli *Client, tun tunnel.Tunnel) *Worker{
	var localHost = defaultHost
	if len(tun.LocalHost) > 0 {
		localHost = tun.LocalHost
	}
	return &Worker{
		cli: cli,
		tun: tun,
		localAddress:  localHost + ":" + strconv.Itoa(int(tun.LocalPort)),
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

func (w *Worker) start() {
	var proxyConn net.Conn
	var err error

	proxyConn, err = w.cli.newConn()
	if err != nil {
		return
	}

	var bridge = transfer.NewBridge(ft, proxyConn, proxyConn)
	_ = bridge.Write(&cmd.RegisterProxy{Tunnel: w.tun, ClientId: w.cli.id})

	for {
		localConn, err := w.newLocalConn()
		if err != nil {
			return
		}
		var end bool
		conn.Combine(localConn, proxyConn, func(alive net.Conn, err error) {
			fmt.Println(alive, localConn, proxyConn)
			if alive == proxyConn {
				w.cli.logger.Warn("The local connection is disconnected:", err)
				_ = localConn.Close()
			} else {
				w.cli.logger.Warn("The proxy connection is disconnected:", err)
				_ = proxyConn.Close()
				end = true
			}
		})
		time.Sleep(time.Duration(2)*time.Hour)
		if end {
			break
		}
	}

	w.cli.logger.Info("The worker is closed")
}