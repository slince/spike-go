package server

import (
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/protol"
	"net"
)

// When receive message.
func OnServerMessage(ev *event.Event) {

}

func RegisterLogListeners(dispatcher *event.Dispatcher) {
	// 注册系统运行
	dispatcher.On(EventServerInit, event.NewListener(OnServerMessage))
	// 注册收到错误消息
	dispatcher.On(EventUnknownMessage, event.NewListener(OnServerMessage))
}


// 当收到消息时
func OnMessage(ev *event.Event){
	// 消息
	message, _ := ev.Parameters["message"]
	msg, _ := message.(protol.Protocol)
	// server
	server, _ := ev.Parameters["server"]
	ser, _ := server.(*Server)
	// connection
	connection, _ := ev.Parameters["connection"]
	conn, _ := connection.(net.Conn)

	messageFactory := MessageHandlerFactory{
		Server: ser,
		Conn: conn,
	}

	var handler MessageHandler

	switch msg.Action {
	case "register":
		handler = messageFactory.NewAuthHandler()
	case "register_tunnel":
		handler = messageFactory.NewRegisterTunnelHandler()
	case "register_proxy":
		handler = messageFactory.NewRegisterProxyHandler()
	}
	ev.Parameters["handler"] = handler
}

//注册监听者
func RegisterSystemListener(dispatcher *event.Dispatcher) {
	//注册收到消息时的事件
	dispatcher.On(EventMessage, event.NewListener(OnMessage))
}
