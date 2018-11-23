package server

import (
	"fmt"
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/log"
	"github.com/slince/spike-go/protol"
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
	Authentication Authentication
	// 客户端
	Clients map[string]*Client
	// Logger
	Logger *log.Logger
	//logFile
	logFile string
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
	server.Logger.Info("The server is running...")
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
func (server *Server) SendMessage(connection net.Conn, message *protol.Protocol){
	jsonString := message.ToString()
	connection.Write([]byte(jsonString))
}

// Send message to client
func (server *Server) SendMessageToClient(client *Client, message *protol.Protocol) {
	server.SendMessage(client.Connection, message)
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
func (server *Server) FindChunkServer(id string) ChunkServer{
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
	RegisterSystemListener(server.Dispatcher)
	// 注册日志监听者
	RegisterLogListeners(server.Dispatcher)
}

// handle connection from client.
func (server *Server) handleConnection(conn net.Conn) {
	server.Logger.Info("Accepted a connection.")

	// 预读多条message
	rd := protol.NewReader(conn)
	messages, err := rd.Read()

	if err != nil {
		server.Logger.Error(err)
		conn.Close()
	}

	for _, message := range messages {
		server.Logger.Info("Received a message:\r\n" + message.ToString())
		err := server.handleMessage(message, conn)
		if err != nil {
			server.Logger.Error(err) //消息处理失败直接关闭
			conn.Close()
		}
	}
}

// Handle message
func (server *Server) handleMessage(message *protol.Protocol, conn net.Conn) error {
	// fire event
	ev := event.NewEvent("message", map[string]interface{}{
		"message":  message,
		"server": server,
		"connection": conn,
	})
	server.Dispatcher.Fire(ev)

	// 获取handler
	hd, ok := ev.Parameters["handler"]
	// 收到不知名的报文
	if !ok {
		ev = event.NewEvent("unknownMessage", map[string]interface{}{"message":  message})
		server.Dispatcher.Fire(ev)
		server.Logger.Warn("receive a unknown message")
		return fmt.Errorf("receive a unknown message")
	}
	// 处理消息
	err := hd.(MessageHandler).Handle(message)
	// 有处理错误直接关闭
	if err != nil {
		server.Logger.Warn("error when handle message")
		conn.Write([]byte(err.Error()))
		conn.Close()
	}
	return nil
}


// Creates a new server.
func NewServer(address string, logFile string) *Server {
	logger := log.NewLogger()
	logger.SetLogFile(logFile).EnableConsole() //开启文件日志和控制台日志
	return &Server{
		address,
		nil,
		event.NewDispatcher(),
		NewSimplePasswordAuth("spike", "spike"),
		make(map[string]*Client, 0),
		logger,
		logFile,
	}
}