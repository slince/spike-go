package server

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/transfer"
)

var ft = transfer.NewFactory()
var types = map[transfer.MsgType]transfer.Command{
	cmd.TypePing:     &cmd.ClientPing{},
	cmd.TypeLogin:    &cmd.Login{},
	cmd.TypeRegisterTunnel: &cmd.RegisterTunnel{},
	cmd.TypeRegisterProxy: &cmd.RegisterProxy{},
}

func init(){
	ft.RegisterTypes(types)
}

