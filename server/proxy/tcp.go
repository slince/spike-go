package proxy

import (
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"os"
	"strconv"
	"time"
)

type TcpHandler struct {
	logger        *log.Logger
	proxyConnPool *conn.Pool
	tun tunnel.Tunnel

	listener      net.Listener
	listenAddress string
	localAddress string
	handleConnCallback func(pubConn net.Conn)
	stop chan bool
}

func NewTcpHandler(logger *log.Logger, connPool *conn.Pool, tun tunnel.Tunnel) *TcpHandler{
	var handler = &TcpHandler{
		logger: logger,
		proxyConnPool: connPool,
		tun: tun,
		listenAddress: net.JoinHostPort("0.0.0.0", strconv.Itoa(tun.ServerPort)),
		localAddress: net.JoinHostPort(tun.LocalHost, strconv.Itoa(tun.LocalPort)),
		stop: make(chan bool, 1),
	}
	handler.handleConnCallback = handler.handleConn
	return handler
}

func (tcp *TcpHandler) Listen() (chan bool, error){
	var address, err = net.ResolveTCPAddr("tcp", tcp.listenAddress)
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return nil, err
	}
	tcp.listener = listener
	var stop = make(chan bool, 1)
	go func() {
		Handle:
		for {
			select {
			case <- tcp.stop:
				break Handle
			default:
				_ = listener.SetDeadline(time.Now().Add(time.Second * 5))
				var con, err1 = listener.Accept()
				if err1 != nil {
					if os.IsTimeout(err1) {
						break
					}
					err = err1
					break Handle
				}
				go tcp.handleConnCallback(con)
			}
		}
		stop <- true
	}()
	return stop, nil
}

func (tcp *TcpHandler) AddProxyConn(proxyConn net.Conn) {
	tcp.proxyConnPool.Put(proxyConn)
}

func (tcp *TcpHandler) Close() {
	tcp.stop <- true
	//_ = tcp.listener.Close()
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
