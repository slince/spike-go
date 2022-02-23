package client

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/transfer"
)

var ft = transfer.NewFactory()

var types = map[transfer.MsgType]transfer.Command{
	cmd.TypePong:     &cmd.ServerPong{},
	cmd.TypeLoginRes: &cmd.LoginRes{},
	cmd.TypeRegisterTunnelRes: &cmd.RegisterTunnelRes{},
	cmd.TypeRequestProxy: &cmd.RequestProxy{},
}

func init(){
	ft.RegisterTypes(types)
}
