package server

import (
	"github.com/slince/jinbox/server/chunk_server"
	"net"
)

type Client struct{
	Id string `json:"id"`
	Connection net.Conn `json:"-"`
	ChunkServers []chunk_server.ChunkServer `json:"-"`
}
