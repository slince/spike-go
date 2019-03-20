package server

import (
	"errors"
	"fmt"
	"github.com/rs/xid"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
)

// 消息处理器接口
type MessageHandler interface {
	// Handle the message
	Handle(message *protol.Protocol) error
}

type Handler struct{
	conn   net.Conn
	server *Server
}

// 客户端注册时消息处理器
type AuthHandler struct{
	Handler
}

func (hdl *AuthHandler) Handle(message *protol.Protocol) error{
	//验证客户端凭证
	auth, ok := message.Body["auth"]
	if !ok {
		msg := &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "403"},
			Body: map[string]interface{}{"error": "bad request"},
		}
		protol.WriteMsg(hdl.conn, msg)
		return errors.New("bad request")
	}
	err := hdl.server.Authentication.Auth(auth.(map[string]interface{}))
	var msg *protol.Protocol
	if err != nil {
		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "403"},
		}
	} else {
		// create one new client
		client := newClient(hdl.conn)
		hdl.server.Clients[client.Id] = client
		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "200"},
			Body: map[string]interface{}{"client": client},
		}
	}
	protol.WriteMsg(hdl.conn, msg)
	return nil
}

// 心跳包处理器
type PingHandler struct{
	Handler
}

func (hdl *PingHandler) Handle(message *protol.Protocol) error{
	msg := &protol.Protocol{
		Action: "pong",
	}
	protol.WriteMsg(hdl.conn, msg)
	return nil
}


// 需要验证之后的消息处理器
type RequireAuthHandler struct {
	Handler
	client *Client
}

func (hdl *RequireAuthHandler) isAuthenticated(message *protol.Protocol) bool{
	clientId, ok := message.Headers["client-id"]
	if !ok {
		return false
	}
	if client, ok := hdl.server.Clients[clientId]; ok{
		hdl.client = client
		return true
	}
	return false
}

// 客户端注册隧道时的消息处理器
type RegisterTunnelHandler struct{
	RequireAuthHandler
}

func (hdl *RegisterTunnelHandler) Handle(message *protol.Protocol) error{
	if !hdl.isAuthenticated(message) {
		return errors.New("the client is not authorized")
	}

	tunnelsInfo, ok := message.Body["tunnels"]
	if !ok {
		return errors.New("missing tunnel info")
	}
	infos := tunnelsInfo.([]interface{})
	var details = make([]map[string]interface{}, len(infos))
	for idx,info := range infos{
		details[idx] = info.(map[string]interface{})
	}

	//创建tunnel
	tunnels := tunnel.NewManyTunnels(details)
	regTunns := make([]tunnel.Tunnel, 0)
	var chunkServers = make([]ChunkServer, 0)
	for _,tun := range tunnels {
		//如果tunnel已经注册则拒绝再次注册
		if hdl.server.checkTunExists(tun) {
			msg := &protol.Protocol{
				Action: "register_tunnel_response",
				Headers: map[string]string{"code": "1"},
				Body: map[string]interface{}{
					"error": "The tunnel has been registered",
					"tunnel": tun,
				},
			}
			protol.WriteMsg(hdl.conn, msg)
			continue
		}
		//创建对应的chunk server
		chunkServer,err := newChunkServer(tun, hdl.server, hdl.client)
		if err != nil {
			msg := &protol.Protocol{
				Action: "register_tunnel_response",
				Headers: map[string]string{"code": "2"},
				Body: map[string]interface{}{
					"error": "Error create chunk server.",
					"tunnel": tun,
				},
			}
			hdl.server.Logger.Warn("fail to create chunk server for the tunnel", err)
			protol.WriteMsg(hdl.conn, msg)
			continue
		}
		regTunns = append(regTunns, tun)
		chunkServers = append(chunkServers, chunkServer)
		// start chunk server
		go chunkServer.run()
	}
	//如果有成功
	if len(regTunns) > 0 {
		//追加入客户端的chunk servers,
		hdl.client.chunkServers = append(hdl.client.chunkServers, chunkServers...)
		//注册成功的客户端
		msg := &protol.Protocol{
			Action: "register_tunnel_response",
			Headers: map[string]string{"code": "200"},
			Body: map[string]interface{}{"tunnels": regTunns},
		}

		protol.WriteMsg(hdl.conn, msg)
		return nil
	} else {
		return errors.New("no tunnel is registered")
	}
}

// 创建chunk server
func newChunkServer(tn tunnel.Tunnel, server *Server, client *Client) (ChunkServer,error){
	var chunkServer ChunkServer
	//生成tunnel的id
	tunnelId := xid.New().String()

	switch tn := tn.(type) {
	case *tunnel.TcpTunnel:
		tn.Id = tunnelId
		chunkServer = &TcpChunkServer{
			tunnel: tn,
			client: client,
			server: server,
			pubConns: make(map[string]*PublicConn, 0),
			pubConnsChan: make(chan *PublicConn, 0),
			closeChan: make(chan int, 1),
		}
	case *tunnel.HttpTunnel:
		tn.Id = tunnelId
		tcpChunkServer := TcpChunkServer{
			tunnel: &tn.TcpTunnel,
			client: client,
			server: server,
			pubConns: make(map[string]*PublicConn, 0),
			pubConnsChan: make(chan *PublicConn, 0),
			closeChan: make(chan int, 1),
		}
		chunkServer = &HttpChunkServer{
			TcpChunkServer: tcpChunkServer,
		}
	default:
		return nil, fmt.Errorf("bad tunnel")
	}
	return chunkServer,nil
}

// 注册代理消息处理器
type RegisterProxyHandler struct{
	RequireAuthHandler
}

func (hdl *RegisterProxyHandler) Handle(message *protol.Protocol) error{
	if !hdl.isAuthenticated(message) {
		return errors.New("the client is not authorized")
	}

	tunId, ok := message.Headers["tunnel-id"]
	if !ok {
		return fmt.Errorf("missing tunnel id")
	}
	chunkServer := hdl.server.findChunkServerByTunId(tunId)
	if chunkServer == nil {
		return fmt.Errorf("the chunk server %s is not found", tunId)
	}
	pubConnId,ok := message.Headers["pub-conn-id"]
	if !ok { //错误的注册代理协议
		return fmt.Errorf("missing public id")
	}
	// set proxy conn
	chunkServer.setProxyConn(pubConnId, hdl.conn)
	return nil
}

// 消息处理器创建工厂
type MessageHandlerFactory struct {
	Conn net.Conn
	Server *Server
}

func (factory MessageHandlerFactory) newHandler() Handler{
	return Handler{
		factory.Conn,
		factory.Server,
	}
}

func (factory MessageHandlerFactory) NewAuthHandler() MessageHandler{
	var handler MessageHandler
	handler = &AuthHandler{
		factory.newHandler(),
	}
	return handler
}

func (factory MessageHandlerFactory) NewPingHandler() MessageHandler{
	var handler MessageHandler
	handler = &PingHandler{
		factory.newHandler(),
	}
	return handler
}

func (factory MessageHandlerFactory) NewRegisterTunnelHandler() MessageHandler{
	var handler MessageHandler
	handler = &RegisterTunnelHandler{
		RequireAuthHandler{
			factory.newHandler(),
			nil,
		},
	}
	return handler
}

func (factory MessageHandlerFactory) NewRegisterProxyHandler() MessageHandler{
	var handler MessageHandler
	handler = &RegisterProxyHandler{
		RequireAuthHandler{
			factory.newHandler(),
			nil,
		},
	}
	return handler
}