package server

import (
	"net"
)

type Client struct{
	Id           string        `json:"id"`
	Conn         net.Conn      `json:"-"`
	ChunkServers []ChunkServer `json:"-"`
}

// close the client
func (client *Client) close() {
	for _, cServer := range client.ChunkServers {
		cServer.Close()
	}
	client.Conn.Close()
}