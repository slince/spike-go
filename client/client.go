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
)

type Client struct {
	// client id
	Id string
	// Server Address
	ServerAddress string
	// Logger
	Logger *log.Logger
	// identifier
	Auth map[string]interface{}
	// Tunnels
	tunnels []tunnel.Tunnel
	// 控制连接
	controlConn net.Conn
	// event dispatcher
	Dispatcher *event.Dispatcher
}

// Run client
func (client *Client) Start() {
	conn,err := net.Dial("tcp", client.ServerAddress)
	if err != nil {
		panic(err)
	}
	client.Logger.Info("the client has been connected to the server")
	client.controlConn = conn
	client.handleControlConnection()
}

// Close the client
func (client *Client) Close() {

}

// 处理控制连接
func (client *Client) handleControlConnection() {
	// 第一步获取授权
	client.sendAuthRequest()
	reader := protol.NewReader(client.controlConn)

	for {
		// 监听消息
		messages, err := reader.Read()
		if err != nil {
			client.Logger.Error(err) //忽略读取
		}
		for _, message := range messages {
			client.handleMessage(message)
		}
	}
}

// 处理消息
func (client *Client) handleMessage(message *protol.Protocol) error {
	// fire event
	ev := event.NewEvent("message", map[string]interface{}{
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
	if hdl, ok := hd.(MessageHandler); ok {
		err := hdl.Handle(message)
		// 有处理错误直接关闭
		if err != nil {
			client.controlConn.Write([]byte(err.Error()))
			client.controlConn.Close()
		}

	}
	return nil
}

// find tunnel by id
func (client *Client) findTunnelById(id string) (tunnel.Tunnel, error) {
	for _, tunnel := range client.tunnels {
		if tunnel.GetId() == id {
			return tunnel, nil
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
	client.controlConn.Write([]byte(message.ToString()))
}

func NewClient(ServerAddress string, tunnels []tunnel.Tunnel, LogFile string) *Client {
	logger := log.NewLogger()
	logger.SetLogFile(LogFile)
	return &Client{
		ServerAddress,
		logger,
		tunnels,
		Dispatcher: event.NewDispatcher(),
	}
}