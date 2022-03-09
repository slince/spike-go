package cmd

import (
	"github.com/slince/spike/pkg/transfer"
	"net"
)

var ft = transfer.NewFactory()
var types = map[transfer.MsgType]transfer.Command{
	TypePong:     &ServerPong{},
	TypeLoginRes: &LoginRes{},
	TypeRegisterTunnelRes: &RegisterTunnelRes{},
	TypeRequestProxy: &RequestProxy{},

	TypePing:     &ClientPing{},
	TypeLogin:    &Login{},
	TypeRegisterTunnel: &RegisterTunnel{},
	TypeRegisterProxy: &RegisterProxy{},

	TypeViewProxy: &ViewProxy{},
	TypeViewProxyResp: &ViewProxyResp{},
}

func init(){
	ft.RegisterTypes(types)
}

func NewBridge(conn net.Conn) *transfer.Bridge{
	return transfer.NewBridge(ft, conn, conn)
}