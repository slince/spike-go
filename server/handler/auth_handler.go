package handler

import (
	"github.com/rs/xid"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/server"
)

// 客户端注册时消息处理器
type AuthHandler struct{
	Handler
}

func (hd *AuthHandler) Handle(message *protol.Protocol){
	//验证客户端凭证
	err := hd.server.Authentication.Auth(message.Body)

	var msg *protol.Protocol
	if err != nil {
		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "403"},
		}
	} else {
		guid := xid.New().String()
		client := &server.Client{
			Connection: hd.connection,
			Id: guid,
		}
		hd.server.Clients[guid] = client

		msg = &protol.Protocol{
			Action: "auth_response",
			Headers: map[string]string{"code": "200"},
			Body: map[string]interface{}{"client": client},
		}
	}
	hd.server.SendMessage(hd.connection, msg)
}

