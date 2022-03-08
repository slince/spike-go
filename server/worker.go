package server

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
)

type Worker struct {
	ser *Server
	tun        tunnel.Tunnel
	conn       net.Conn
	bridge     *transfer.Bridge
	cli *Client
	socket     net.Listener
	proxyConns *conn.Pool
}

func newWorker(ser *Server, tun tunnel.Tunnel, conn net.Conn, bridge *transfer.Bridge, cli *Client) *Worker {
	var worker = &Worker{
		ser, tun, conn, bridge, cli, nil, nil,
	}
	worker.Init()
	return worker
}

func (w *Worker) Init() {
	w.proxyConns = conn.NewPool(10, func(pool *conn.Pool) {
		w.ser.logger.Info("Request to client for proxy connection")
		err := w.requestProxy()
		if err != nil {
			w.ser.logger.Error("Failed to send request proxy command")
		}
	})
}

func (w *Worker) Start() (err error) {
	var address = "0.0.0.0:" + strconv.Itoa(w.tun.ServerPort)
	socket, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	w.socket = socket

	go func() {
		for {
			var con, err1 = socket.Accept()
			if err1 != nil {
				err = err1
				return
			}
			go w.handleConn(con)
		}
	}()
	return
}

func (w *Worker) Close() {
	_ = w.socket.Close()
}

func (w *Worker) addProxyConn(conn net.Conn) {
	w.proxyConns.Put(conn)
}

func (w *Worker) requestProxy() error {
	return w.bridge.Write(&cmd.RequestProxy{ServerPort: w.tun.ServerPort})
}

func (w *Worker) handleConn(pubConn net.Conn) {
	w.ser.logger.Trace("Accept a public connection:", pubConn.RemoteAddr().String())
	var proxyConn = w.proxyConns.Get()
	conn.Combine(proxyConn, pubConn)
}
