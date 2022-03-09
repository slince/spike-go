package proxy

import "net"

type Handler interface {
	Listen(serverPort int) error
	Close()
	AddProxyConn(proxyConn net.Conn)
}

