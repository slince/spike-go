package server

import (
	"github.com/rs/xid"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type Client struct{
	Id           string        `json:"id"`
	controlConn  net.Conn      `json:"-"`
	chunkServers []ChunkServer `json:"-"`
	tunnels map[string]tunnel.Tunnel
}

// close the client
func (client *Client) close() {
	for _, chunkServer := range client.chunkServers {
		chunkServer.close()
	}
	// 关闭当前控制连接
	client.controlConn.Close()
}

// create one client
func newClient(controlConn net.Conn) *Client{
	return &Client{
		Id: xid.New().String(),
		controlConn: controlConn,
		tunnels: make(map[string]tunnel.Tunnel, 0),
	}
}