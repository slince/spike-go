package proxy

import (
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"net"
	"strconv"
)

type TcpHandler struct {
	logger        *log.Logger
	listener      net.Listener
	localAddress string
	proxyConnPool *conn.Pool
	handleConnCallback func(pubConn net.Conn)
}

func NewTcpHandler(logger *log.Logger, connPool *conn.Pool, localAddress string) *TcpHandler{
	var handler = &TcpHandler{
		logger: logger,
		proxyConnPool: connPool,
		localAddress: localAddress,
	}
	handler.handleConnCallback = handler.handleConn
	return handler
}

func (tcp *TcpHandler) Listen(serverPort int) error{
	var address, err = net.ResolveTCPAddr("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(serverPort)))
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}
	tcp.listener = listener

	go func() {
		for {
			var con, err1 = listener.Accept()
			if err1 != nil {
				err = err1
				break
			}

			go tcp.handleConnCallback(con)
		}
	}()
	return nil
}

func (tcp *TcpHandler) AddProxyConn(proxyConn net.Conn) {
	tcp.proxyConnPool.Put(proxyConn)
}

func (tcp *TcpHandler) Close() {
	_ = tcp.listener.Close()
}

func (tcp *TcpHandler) handleConn(pubConn net.Conn) {
	tcp.logger.Trace("Accept a public connection:", pubConn.RemoteAddr().String())
	var proxyConn, err = tcp.proxyConnPool.Get()
	if err != nil {
		tcp.logger.Error("Failed to get proxy conn from client, error", err)
		pubConn.Close()
		return
	}
	conn.Combine(proxyConn, pubConn)
}
