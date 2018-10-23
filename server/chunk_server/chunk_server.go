package chunk_server

import "github.com/slince/jinbox/tunnel"

type ChunkServer interface {
	//启动监听服务
	Run()
	//获取对应的tunnel
	GetTunnel() tunnel.Tunnel
}

// 监听公网接口
type TcpChunkServer struct {
	Tunnel tunnel.Tunnel
}

func (chunkServer *TcpChunkServer) Run() {

}

func (chunkServer *TcpChunkServer) GetTunnel() tunnel.Tunnel{
	return chunkServer.Tunnel
}

// http chunk server
type HttpChunkServer struct {
	TcpChunkServer
}