package chunk_server

import (
	"github.com/rs/xid"
	"net"
)

type PublicConn struct {
	Id string
	Conn net.Conn
	ProxyConnChan *chan net.Conn
}

// Create a public connection.
func NewPublicConn(conn net.Conn) *PublicConn {
	ch := make(chan net.Conn)
	return &PublicConn{
		Id: xid.New().String(),
		Conn: conn,
		ProxyConnChan: &ch,
	}
}