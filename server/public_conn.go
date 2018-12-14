package server

import (
	"github.com/rs/xid"
	"io"
	"net"
	"sync"
)

type PublicConn struct {
	Id string
	Conn net.Conn
	ProxyConnChan chan net.Conn
	pubLock sync.RWMutex
	proxyLock sync.RWMutex
}

func (pubConn *PublicConn) Pipe(conn net.Conn) {

	defer conn.Close()
	defer pubConn.Conn.Close()

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		for {
			pubConn.pubLock.RLock()
			pubConn.proxyLock.Lock()
			io.Copy(conn, pubConn.Conn)
			pubConn.pubLock.RUnlock()
			pubConn.proxyLock.RLock()
		}
	}()
	go func() {
		for {
			pubConn.pubLock.Lock()
			pubConn.proxyLock.RLock()
			io.Copy(pubConn.Conn, conn)
			pubConn.pubLock.Unlock()
			pubConn.proxyLock.RUnlock()
		}
	}()

	wait.Wait()
}

// Create a public connection.
func NewPublicConn(conn net.Conn) *PublicConn {
	ch := make(chan net.Conn)
	return &PublicConn{
		Id: xid.New().String(),
		Conn: conn,
		ProxyConnChan: ch,
	}
}