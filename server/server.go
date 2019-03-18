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
	//事件处理器
	Dispatcher *event.Dispatcher
	//验证对象
	Authentication Authentication
	// 客户端
	Clients map[string]*Client
	// Logger
	Logger *log.Logger
	// control conn control
	controlConnChan chan net.Conn
	// closed
	closeChan chan int
}

// Close the server
func (server *Server) Stop(){
	server.closeChan <- 1
}

// Run the server
func (server *Server) Run() {
	// register all listeners
	server.registerListeners()
	// 监听端口
	var err error
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		panic(err.Error())
	}
	server.Logger.Info("The server is running...")

	go server.acceptControlConn(listener) // 接收请求
	go server.processControlConns()        // 消费所有控制请求
	// check if closed.
	<-server.closeChan
}

// 接收请求
func (server *Server) acceptControlConn(listener net.Listener) {
	for {
		select {
		case <- server.closeChan:
			listener.Close()
		default:
			conn, err := listener.Accept()
			if err != nil {
				// handle error
				continue
			}
			server.controlConnChan <- conn
		}
	}
}

// process conns
func (server *Server) processControlConns() {
	for {
		select {
		case <- server.closeChan:
			break
		default:
			conn := <- server.controlConnChan
			go server.handleControlConn(conn)
		}
	}
}

// handle connection from client.
func (server *Server) handleControlConn(conn net.Conn) {
	server.Logger.Info("Accepted a connection.")
	// 预读多条message
	reader := protol.NewReader(conn)
	for {
		messages, err := reader.Read()
		if err != nil { //如果读取失败则关闭次连接
			server.Logger.Error(err)
			client := server.findClientByConn(conn)
			if client != nil {
				client.close()
			}
			break
		}
		for  _, message := range messages {
			server.Logger.Info("Received a message:\r\n" + message.ToString())
			server.handleMessage(message, conn)
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
	hdl, ok := ev.Parameters["handler"]
	// 收到不知名的报文
	if !ok {
		ev = event.NewEvent("unknownMessage", map[string]interface{}{"message":  message})
		server.Dispatcher.Fire(ev)
		server.Logger.Warn("receive a unknown message")
		return fmt.Errorf("receive a unknown message")
	}
	// 处理消息
	err := hdl.(MessageHandler).Handle(message)
	// 有处理错误直接关闭
	if err != nil {
		server.Logger.Warn(err)
		conn.Close()
	}
	return nil
}

// Send message to connection
func (server *Server) sendMessage(connection net.Conn, message *protol.Protocol){
	jsonString := message.ToString()
	connection.Write([]byte(jsonString))
}

// Send message to client
func (server *Server) sendToClient(client *Client, message *protol.Protocol) {
	server.sendMessage(client.controlConn, message)
}

// Checks whether tunnel is registered.
func (server *Server) checkTunExists(tunnel tunnel.Tunnel) bool {
	for _,client := range server.Clients {
		for _, tun := range client.tunnels {
			if tun.MatchTunnel(tunnel) {
				return true
			}
		}
	}
	return false
}

// find chunk server by its tunnel
func (server *Server) findChunkServerByTunId(id string) ChunkServer{
	for _, client := range server.Clients {
		for _, chunkServer := range client.chunkServers {
			if chunkServer.getTunnel().GetId() == id {
				return chunkServer
			}
		}
	}
	return nil
}

// find client
func (server *Server) findClientByConn(conn net.Conn) *Client{
	for _,client := range server.Clients {
		if client.controlConn == conn {
			return client
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

// Creates a new server.
func NewServer(configuration *Configuration) *Server {
	logger := log.NewLogger()
	logger.SetLogFile(configuration.Log["file"]).EnableConsole() //开启文件日志和控制台日志

	authentication := NewSimplePasswordAuth(
		configuration.Auth["username"],
		configuration.Auth["password"],
	)
	return &Server{
		Address:         configuration.Address,
		Dispatcher:      event.NewDispatcher(),
		Authentication:  authentication,
		Clients:         make(map[string]*Client, 0),
		Logger:          logger,
		controlConnChan: make(chan net.Conn),
		closeChan: make(chan int, 1),
	}
}