package chunk_server

import "github.com/slince/jinbox/tunnel"

type ChunkServer interface {
	Run()
}

// 监听公网接口
type TcpChunkServer struct {
	Tunnel tunnel.Tunnel
}

func (chunkServer *TcpChunkServer) Run() {

}

// http chunk server
type HttpChunkServer struct {
	TcpChunkServer
}