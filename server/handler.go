package server

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/server/chunk_server"
	"github.com/slince/spike-go/tunnel"
	"net"
)
// 消息处理器接口
type MessageHandler interface {
	// Handle the message
	Handle(message *protol.Protocol) error
}

type Handler struct{
	connection net.Conn
	server *Server
}

// 客户端注册时消息处理器
type AuthHandler struct{
	Handler
}

func (hd *AuthHandler) Handle(message *protol.Protocol){
	//验证客户端凭证
	err := hd.server.Authentication.Auth(message.Body)

	var msg *protol.Protocol
	if err != nil {
		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "403"},
		}
	} else {
		guid := xid.New().String()
		client := &Client{
			Connection: hd.connection,
			Id: guid,
		}
		hd.server.Clients[guid] = client

		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "200"},
			Body: map[string]interface{}{"client": client},
		}
	}
	hd.server.SendMessage(hd.connection, msg)
}

// 心跳包处理器
type PingHandler struct{
	Handler
}

func (hd *PingHandler) Handle(message *protol.Protocol){
	msg := &protol.Protocol{
		Action: "pong",
	}
	hd.server.SendMessage(hd.connection, msg)
}


// 需要验证之后的消息处理器
type RequireAuthHandler struct {
	Handler
	client *Client
}

func (hd *RequireAuthHandler) isAuthenticated(message *protol.Protocol) bool{
	clientId, ok := message.Headers["client-id"]
	if !ok {
		return false
	}
	if client, ok := hd.server.Clients[clientId]; ok{
		hd.client = client
		return true
	}
	return false
}

// 注册代理消息处理器
type RegisterProxyHandler struct{
	RequireAuthHandler
}

func (hd *RegisterProxyHandler) Handle(message *protol.Protocol){
	tunnelId, ok := message.Headers["tunnel-id"]
	if !ok {
		hd.connection.Write([]byte("missing tunnel id"))
		hd.connection.Close()
		return
	}
	chunkServer := hd.server.FindChunkServer(tunnelId)
	if chunkServer == nil {
		hd.connection.Write([]byte(fmt.Sprintf("the chunk server %s is not found", tunnelId)))
		hd.connection.Close()
		return
	}
	publicConnectionId,ok := message.Headers["public-connection-id"]
	if !ok { //错误的注册代理协议
		hd.connection.Write([]byte("missing public id"))
		hd.connection.Close()
	}
	// set proxy connection
	chunkServer.SetProxyConnection(publicConnectionId, hd.connection)
}

// 客户端注册隧道时的消息处理器
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
			Tunnel: &tn.TcpTunnel,
		}
		chunkServer = &chunk_server.HttpChunkServer{
			TcpChunkServer: tcpChunkServer,
		}
	default:
		return nil, fmt.Errorf("bad tunnel")
	}
	return chunkServer,nil
}
