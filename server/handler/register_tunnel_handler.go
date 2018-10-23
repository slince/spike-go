package handler

import (
	"fmt"
	"github.com/rs/xid"
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
	registeredTunnels := make([]tunnel.Tunnel, 0)

	var chunkServers = make([]chunk_server.ChunkServer, len(tunnels))
	for _,tn := range tunnels {
		//如果tunnel已经注册则拒绝再次注册
		if hd.server.IsTunnelRegistered(tn) {
			msg := &protol.Protocol{
				Action: "register_tunnel_response",
				Headers: map[string]string{"code": "1"},
				Body: map[string]interface{}{
					"error": "The tunnel has been registered",
					"tunnel": tn,
				},
			}
			go hd.server.SendMessage(hd.connection, msg)
			continue
		}
		//创建对应的chunk server
		chunkServer,err := newChunkServer(tn)
		if err != nil {
			msg := &protol.Protocol{
				Action: "register_tunnel_response",
				Headers: map[string]string{"code": "2"},
				Body: map[string]interface{}{
					"error": "Error create chunk server.",
					"tunnel": tn,
				},
			}
			go hd.server.SendMessage(hd.connection, msg)
			continue
		}
		registeredTunnels = append(registeredTunnels, tn)
		chunkServers = append(chunkServers, chunkServer)
	}
	//追加入客户端的chunk servers
	hd.client.ChunkServers = append(hd.client.ChunkServers, chunkServers...)
	//注册成功的客户端
	msg := &protol.Protocol{
		Action: "register_tunnel_response",
		Headers: map[string]string{"code": "200"},
		Body: map[string]interface{}{"tunnels": registeredTunnels},
	}
	go hd.server.SendMessage(hd.connection, msg)
}

// 创建chunk server
func newChunkServer(tn tunnel.Tunnel) (chunk_server.ChunkServer,error){
	var chunkServer chunk_server.ChunkServer
	//生成tunnel的id
	tunnelId := xid.New().String()
	switch tn := tn.(type) {
	case *tunnel.TcpTunnel:
		tn.Id = tunnelId
		chunkServer = &chunk_server.TcpChunkServer{
			Tunnel: tn,
		}
	case *tunnel.HttpTunnel:
		tn.Id = tunnelId
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

