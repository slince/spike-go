package handler

import (
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
	var chunkServers = make([]*chunk_server.ChunkServer, len(tunnels))
	for _,tn := range tunnels {
		chunkServer := newChunkServer(tn)
		chunkServers = append(chunkServers, chunkServer)
	}
	//追加进入chunk管理
	clientChunkServers,ok := hd.server.ChunkServers[hd.client.Id]
	if !ok {
		clientChunkServers = make([]*chunk_server.ChunkServer, len(chunkServers))
		hd.server.ChunkServers[hd.client.Id] = clientChunkServers
	}
	clientChunkServers = append(clientChunkServers, chunkServers...)
}

func newChunkServer(tn tunnel.Tunnel) *chunk_server.ChunkServer{

}

