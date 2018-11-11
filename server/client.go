package server

import (
	"net"
)

type Client struct{
	Id string `json:"id"`
	Connection net.Conn `json:"-"`
	ChunkServers []ChunkServer `json:"-"`
}
