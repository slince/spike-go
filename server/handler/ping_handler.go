package handler

import (
	"github.com/slince/jinbox/protol"
)

// 心跳包处理器
type PingHandler struct{
	Handler
}

func (hd *PingHandler) Handle(message *protol.Protocol){
	msg := &protol.Protocol{
		Action: "pong",
	}
	hd.server.SendMessage(hd.connection, msg)
}
