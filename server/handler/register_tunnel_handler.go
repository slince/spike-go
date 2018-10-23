package handler

import (
	"fmt"
	"github.com/slince/jinbox/protol"
	"github.com/slince/jinbox/server/chunk_server"
	"github.com/slince/jinbox/tunnel"
)

// 客户端注册时消息处理器
type RegisterTunnelHandler struct{
	RequireAuthHandler
}

func (hd *RegisterTunnelHandler) Handle(message *protol.Protocol){
	tunnelsInfo, ok := message.Body["tunnels"]
	if !ok {
		return
	}
	tunnelsInfoValue, ok := tunnelsInfo.([]map[string]string)
	if !ok {
		return
	}
	//创建tunnel
	tunnels := tunnel.NewManyTunnels(tunnelsInfoValue)
	var chunkServers = make([]chunk_server.ChunkServer, len(tunnels))
	for _,tn := range tunnels {
		chunkServer,err := newChunkServer(tn)
		if err != nil {
			continue
		}
		chunkServers = append(chunkServers, chunkServer)
	}
	//追加入客户端的chunk servers
	hd.client.ChunkServers = append(hd.client.ChunkServers, chunkServers...)
}

// 创建chunk server
func newChunkServer(tn tunnel.Tunnel) (chunk_server.ChunkServer,error){
	var chunkServer chunk_server.ChunkServer
	switch tn := tn.(type) {
	case *tunnel.TcpTunnel:
		chunkServer = &chunk_server.TcpChunkServer{
			Tunnel: tn,
		}
	case *tunnel.HttpTunnel:
		tcpChunkServer := chunk_server.TcpChunkServer{
			Tunnel:tn,
		}
		chunkServer = &chunk_server.HttpChunkServer{
			TcpChunkServer: tcpChunkServer,
		}
	default:
		return nil, fmt.Errorf("bad tunnel")
	}
	return chunkServer,nil
}

