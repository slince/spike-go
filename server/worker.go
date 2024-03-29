package server

import (
	"fmt"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"github.com/slince/spike/server/proxy"
	"net"
)

type Worker struct {
	ser *Server
	tun        tunnel.Tunnel
	conn       net.Conn
	bridge     *transfer.Bridge
	cli *Client
	handler proxy.Handler
}

func newWorker(ser *Server, tun tunnel.Tunnel, conn net.Conn, bridge *transfer.Bridge, cli *Client) *Worker {
	var worker = &Worker{
		ser, tun, conn, bridge, cli, nil,
	}
	return worker
}

func (w *Worker) Start() error {
	var handler, err = w.createHandler()
	if err != nil {
		return err
	}
	w.handler = handler
	stop, err := handler.Listen()
	if err != nil {
		return err
	}
	go func() {
		<- stop
		close(stop)
		w.ser.logger.Info(fmt.Sprintf("The worker for %d is closed", w.tun.ServerPort))
	}()
	return nil
}

func (w *Worker) createHandler() (proxy.Handler, error){
	var handler proxy.Handler
	var connPool = conn.NewPool(100, 5, func(pool *conn.Pool) {
		w.ser.logger.Info("Request to client for proxy connection")
		err := w.requestProxy()
		if err != nil {
			w.ser.logger.Error("Failed to send request proxy command")
		}
	})
	var err error
	switch w.tun.Protocol {
	case "tcp":
		handler = proxy.NewTcpHandler(w.ser.logger, connPool, w.tun)
	case "udp":
		handler = proxy.NewUdpHandler(w.ser.logger, connPool, w.tun)
	case "http":
		handler = proxy.NewHttpHandler(w.ser.logger, connPool, w.tun, w.tun.Headers)
	default:
		err = fmt.Errorf("unsupported tunel protocol %s", w.tun.Protocol)
	}
	return handler,err
}

func (w *Worker) Close() {
	w.handler.Close()
}

func (w *Worker) addProxyConn(conn net.Conn) {
	w.handler.AddProxyConn(conn)
}

func (w *Worker) requestProxy() error {
	return w.bridge.Write(&cmd.RequestProxy{ServerPort: w.tun.ServerPort})
}
