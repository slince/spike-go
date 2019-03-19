package server

import (
	"github.com/rs/xid"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type Client struct{
	Id           string        `json:"id"`
	ctrlConn     net.Conn      `json:"-"`
	chunkServers []ChunkServer `json:"-"`
	tunnels      map[string]tunnel.Tunnel
}

// close the client
func (client *Client) close() {
	for _, chunkServer := range client.chunkServers {
		chunkServer.close()
	}
	// 关闭当前控制连接
	client.ctrlConn.Close()
}

// create one client
func newClient(controlConn net.Conn) *Client{
	return &Client{
		Id:       xid.New().String(),
		ctrlConn: controlConn,
		tunnels:  make(map[string]tunnel.Tunnel, 0),
	}
}