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
	msg := ev.Parameters["message"].(*protol.Protocol)
	// server
	ser := ev.Parameters["server"].(*Server)
	// connection
	conn := ev.Parameters["connection"].(net.Conn)
	messageFactory := MessageHandlerFactory{
		Server: ser,
		Conn: conn,
	}
	var handler MessageHandler
	switch msg.Action {
	case "auth":
		handler = messageFactory.NewAuthHandler()
	case "ping":
		handler = messageFactory.NewPingHandler()
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
