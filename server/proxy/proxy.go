package proxy

import "net"

type Handler interface {
	Listen() (chan bool, error)
	Close()
	AddProxyConn(proxyConn net.Conn)
}

