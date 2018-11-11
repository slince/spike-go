package server

import (
	"fmt"
	"github.com/slince/spike-go/auth"
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/server/chunk_server"
	"github.com/slince/spike-go/server/handler"
	"github.com/slince/spike-go/server/listener"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type Server struct {
	//绑定的地址
	Address string
	//控制套接字
	socket net.Listener
	//事件处理器
	Dispatcher *event.Dispatcher
	//验证对象
	Authentication auth.Authentication
	// 客户端
	Clients map[string]*Client
}

// Run the server
func (server *Server) Run() {
	// register all listeners
	server.registerListeners()

	// 监听端口
	var err error
	server.socket ,err = net.Listen("tcp", server.Address)
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := server.socket.Accept()
		if err != nil {
			// handle error
			continue
		}
		go server.handleConnection(conn)
	}
}

// Send message to connection
func (server *Server) SendMessage(connection net.Conn, message *protol.Protocol) error{
	jsonString, err := message.ToString()
	if err != nil {
		return err
	}
	connection.Write([]byte(jsonString))
	return nil
}

// Send message to client
func (server *Server) SendMessageToClient(client *Client, message *protol.Protocol) error {
	return server.SendMessage(client.Connection, message)
}

// Checks whether tunnel is registered.
func (server *Server) IsTunnelRegistered(tunnel tunnel.Tunnel) bool {
	for _,client := range server.Clients {
		for _, chunkServer := range client.ChunkServers {
			if chunkServer.GetTunnel().MatchTunnel(tunnel) {
				return true
			}
		}
	}
	return false
}

// find chunk server by its tunnel
func (server *Server) FindChunkServer(id string) chunk_server.ChunkServer{
	for _, client := range server.Clients {
		for _, chunkServer := range client.ChunkServers {
			if chunkServer.GetTunnel().GetId() == id {
				return chunkServer
			}
		}
	}
	return nil
}

// Register all listeners
func (server *Server) registerListeners() {
	// 注册系统监听者
	listener.RegisterSystemListener(server.Dispatcher)
	// 注册日志监听者
	listener.RegisterLogListeners(server.Dispatcher)
}

// handle connection from client.
func (server *Server) handleConnection(connection net.Conn) error{
	// 预读多条message
	rd := protol.NewReader(connection)
	messages,err := rd.Read()

	if err != nil {
		return err
	}

	for _, message := range messages {
		err := server.handleMessage(message)
		if err != nil {
			return err
		}
	}

	return nil
}

// Handle message
func (server *Server) handleMessage(message *protol.Protocol) error {
	ev := event.NewEvent("message", map[string]interface{}{"message":  message})
	server.Dispatcher.Fire(ev)

	hd, ok := ev.Parameters["handler"]
	// 收到不知名的报文
	if !ok {
		ev = event.NewEvent("unknownMessage", map[string]interface{}{"message":  message})
		server.Dispatcher.Fire(ev)
		return fmt.Errorf("recieve a unknown message")
	}
	// 处理消息
	hdl, ok := hd.(handler.MessageHandler)
	if ok {
		hdl.Handle(message)
	}
	return nil
}


// Creates a new server.
func NewServer(address string) Server {
	return Server{
		address,
		nil,
		event.NewDispatcher(),
		nil,
		make(map[string]*Client, 0),
	}
}
