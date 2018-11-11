package listener

import (
	"github.com/slince/spike-go/event"
	"github.com/slince/spike-go/server"
)

// When receive message.
func OnServerMessage(ev *event.Event) {

}

func RegisterLogListeners(dispatcher *event.Dispatcher) {
	// 注册系统运行
	dispatcher.On(server.ServerInit, event.NewListener(OnServerMessage))
	// 注册收到错误消息
	dispatcher.On(server.UnknownMessage, event.NewListener(OnServerMessage))
}



