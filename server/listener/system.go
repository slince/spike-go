package listener

import (
	"github.com/slince/jinbox/event"
	"github.com/slince/jinbox/protol"
	"github.com/slince/jinbox/server"
	"github.com/slince/jinbox/server/handler"
)

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
	dispatcher.On(server.Message, event.NewListener(OnMessage))
}