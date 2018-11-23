package client

import (
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/protol"
)

// 当收到消息时
func OnMessage(ev *event.Event){
	// 消息
	msg := ev.Parameters["message"].(*protol.Protocol)
	// client
	client := ev.Parameters["client"].(*Client)
	messageFactory := MessageHandlerFactory{
		client: client,
	}
	var handler MessageHandler
	switch msg.Action {
	case "auth_response":
		handler = messageFactory.NewAuthResponseHandler()
	case "register_tunnel_response":
		handler = messageFactory.NewRegisterTunnelResponseHandler()
	case "request_proxy":
		handler = messageFactory.NewRequestProxyHandler()
	}
	ev.Parameters["handler"] = handler
}

//注册监听者
func RegisterSystemListener(dispatcher *event.Dispatcher) {
	//注册收到消息时的事件
	dispatcher.On(EventMessage, event.NewListener(OnMessage))
}

