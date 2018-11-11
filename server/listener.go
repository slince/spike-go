package server

import (
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/protol"
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
	message, ok := ev.Parameters["message"]
	if !ok {
		return
	}
	msg, ok := message.(protol.Protocol)
	if !ok {
		return
	}

	switch msg.Action {
	case "register":
	case "register_tunnel":
	case "register_proxy":
	}
	ev.Parameters["handler"] = handler
}

//注册监听者
func RegisterSystemListener(dispatcher *event.Dispatcher) {
	//注册收到消息时的事件
	dispatcher.On(EventMessage, event.NewListener(OnMessage))
}
