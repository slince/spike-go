package client

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/log"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
	"runtime"
)

// 初始化常量
const (
	Version = "0.0.1"
	EventClientInit = "init"
	EventClientStart = "start"
	EventMessage = "message"
	EventUnknownMessage = "unknownMessage"
)

type Client struct {
	// client id
	Id string
	// Server Address
	ServerAddress string
	// Logger
	Logger *log.Logger
	// identifier
	Auth map[string]string
	// Tunnels
	Tunnels []tunnel.Tunnel
	// 控制连接
	ControlConn net.Conn
	// event dispatcher
	Dispatcher *event.Dispatcher
}

// Run client
func (client *Client) Start() {
	client.registerListeners()

	conn,err := net.Dial("tcp", client.ServerAddress)
	if err != nil {
		panic(err)
	}
	client.Logger.Info("the client has been connected to the server")
	client.ControlConn = conn

	client.handleControlConnection()
}

// Close the client
func (client *Client) Close() {

}

// 发送消息给服务端
func (client *Client) SendMessage(message *protol.Protocol) (int, error){
	if client.Id != "" {
		if message.Headers == nil {
			message.Headers = make(map[string]string, 1)
		}
		message.Headers["client-id"] = client.Id
	}
	return client.ControlConn.Write(message.ToBytes())
}

// Register all listeners
func (client *Client) registerListeners() {
	// 注册系统监听者
	RegisterSystemListener(client.Dispatcher)
}

// 处理控制连接
func (client *Client) handleControlConnection() {
	// 第一步获取授权
	client.sendAuthRequest()
	reader := protol.NewReader(client.ControlConn)

	for {
		// 监听消息
		messages, err := reader.Read()
		if err != nil {
			client.Logger.Error(err) //忽略读取
		}
		for _, message := range messages {
			client.Logger.Info("Received a message:\r\n" + message.ToString())
			client.handleMessage(message)
		}
	}
}

// 处理消息
func (client *Client) handleMessage(message *protol.Protocol) error {
	// fire event
	ev := event.NewEvent(EventMessage, map[string]interface{}{
		"message":  message,
		"client": client,
	})
	client.Dispatcher.Fire(ev)
	// 获取handler
	hd, ok := ev.Parameters["handler"]
	// 收到不知名的报文
	if !ok {
		ev = event.NewEvent("unknownMessage", map[string]interface{}{"message":  message})
		client.Dispatcher.Fire(ev)
		client.Logger.Warn("receive a unknown message")
		return fmt.Errorf("receive a unknown message")
	}
	// 处理消息
	err := hd.(MessageHandler).Handle(message)
	if err != nil {
		client.Logger.Warn("bad response from the server")
	}
	return nil
}

// find tunnel by id
func (client *Client) findTunnelById(id string) (tunnel.Tunnel, error) {
	for _, tn := range client.Tunnels {
		if tn.GetId() == id {
			return tn, nil
		}
	}
	return nil, errors.New("The tunnel is missing with id")
}

// 发送验证信息给服务端
func (client *Client) sendAuthRequest() {
	message := &protol.Protocol{
		Action: "auth",
		Body: map[string]interface{}{
			"os": runtime.GOOS + runtime.GOARCH,
			"version": Version,
			"auth": client.Auth,
		},
	}
	client.ControlConn.Write([]byte(message.ToString()))
}

func NewClient(configuration *Configuration) *Client {
	// set logger
	logger := log.NewLogger()
	logger.SetLogFile(configuration.Log["file"]).EnableConsole()
	tunnels := createTunnelsWithTunnelConfiguration(configuration.Tunnels)
	return &Client{
		Id: "",
		ServerAddress: configuration.ServerAddress,
		Logger: logger,
		Auth: configuration.Auth,
		Tunnels: tunnels,
		Dispatcher: event.NewDispatcher(),
	}
}

// 创建tunnel
func createTunnelsWithTunnelConfiguration(configurations []TunnelConfiguration) []tunnel.Tunnel{
	var details []map[string]interface{}
	for _, config := range configurations {
		details = append(details, map[string]interface{}{
			"protocol": config.Protocol,
			"local_port": config.LocalPort,
			"server_port": config.ServerPort,
			"host": config.Host,
			"proxy_hosts": config.ProxyHosts,
		})
	}
	return tunnel.NewManyTunnels(details)
}