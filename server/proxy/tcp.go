package proxy

import (
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"net"
	"strconv"
)

type TcpHandler struct {
	logger *log.Logger
	socket net.Listener
	proxyConnPool *conn.Pool
}

func NewTcpHandler(logger *log.Logger, connPool *conn.Pool) *TcpHandler{
	return &TcpHandler{
		logger: logger,
		proxyConnPool: connPool,
	}
}

func (tcp *TcpHandler) Listen(serverPort int) error{
	var address, err = net.ResolveTCPAddr("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(serverPort)))
	socket, err := net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}
	tcp.socket = socket

	go func() {
		for {
			var con, err1 = socket.Accept()
			if err1 != nil {
				err = err1
				break
			}

			go tcp.handleConn(con)
		}
	}()
	return nil
}

func (tcp *TcpHandler) AddProxyConn(proxyConn net.Conn) {
	tcp.proxyConnPool.Put(proxyConn)
}

func (tcp *TcpHandler) Close() {
	_ = tcp.socket.Close()
}

func (tcp *TcpHandler) handleConn(pubConn net.Conn) {
	tcp.logger.Trace("Accept a public connection:", pubConn.RemoteAddr().String())
	var proxyConn, err = tcp.proxyConnPool.Get()
	if err != nil {
		tcp.logger.Error("Failed to get proxy conn from client, error", err)
		pubConn.Close()
	}
	conn.Combine(proxyConn, pubConn)
}
