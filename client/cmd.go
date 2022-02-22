package client

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/transfer"
)

var ft = transfer.NewFactory()

var types = map[transfer.MsgType]transfer.Command{
	cmd.TypePing:     cmd.ClientPing{},
	cmd.TypePong:     cmd.ServerPong{},
	cmd.TypeLogin:    cmd.Login{},
	cmd.TypeLoginRes: cmd.LoginRes{},
}

func init(){
	ft.RegisterTypes(types)
}
