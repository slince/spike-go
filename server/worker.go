package server

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"io"
	"net"
	"strconv"
)

type Worker struct {
	tun        tunnel.Tunnel
	conn       net.Conn
	socket     net.Listener
	bridge     *transfer.Bridge
	proxyConns *conn.Pool
}

func newWorker(tun tunnel.Tunnel, conn net.Conn, bridge *transfer.Bridge) *Worker {
	var worker = &Worker{
		tun, conn, nil, bridge, nil,
	}
	worker.Init()
	return worker
}

func (w *Worker) Init() {
	w.proxyConns = conn.NewPool(10, func() {
		err := w.requestProxy()
		if err != nil {
			logger.Error("Failed to send request proxy command")
		}
	})
}

func (w *Worker) Start() (err error) {
	address := "0.0.0.0:" + strconv.Itoa(int(w.tun.ServerPort))
	socket, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	w.socket = socket
	for {
		var con, err1 = socket.Accept()
		if err1 != nil {
			err = err1
			return
		}
		go w.handleConn(con)
	}
}

func (w *Worker) Close() error {
	return w.socket.Close()
}

func (w *Worker) AddProxyConn(conn net.Conn) {
	w.proxyConns.Put(conn)
}

func (w *Worker) requestProxy() error {
	return w.bridge.Write(cmd.RequestProxy{Tunnel: w.tun})
}

func (w *Worker) handleConn(con net.Conn) {
	var proxyConn = w.proxyConns.Get()
	var readErr = func(src io.Reader) {
		con = src.(net.Conn)
		if src != proxyConn {
			w.proxyConns.Put(con)
		}
		con.Close()
	}
	var writeErr = func(dst io.Writer) {
		con = dst.(net.Conn)
		if dst != proxyConn {
			w.proxyConns.Put(con)
		}
		con.Close()
	}
	conn.Combine(proxyConn, con, readErr, writeErr)
}
