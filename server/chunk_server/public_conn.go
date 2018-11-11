package chunk_server

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
}

func (pubConn *PublicConn) Pipe(conn net.Conn) {

	defer conn.Close()
	defer pubConn.Conn.Close()

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		io.Copy(conn, pubConn.Conn)
		wait.Done()
	}()
	go func() {
		io.Copy(pubConn.Conn, conn)
		wait.Done()
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