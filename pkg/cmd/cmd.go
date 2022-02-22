package cmd

import (
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
)

const (
	TypePing = iota
	TypePong
	TypeLogin
	TypeLoginRes
)

var p = transfer.NewParser()

// ClientPing 客户端pin消息
type ClientPing struct {
	ClientId string `json:"client_id"`
}

func (c *ClientPing) GetType() transfer.MsgType {
	return TypePing
}

// ServerPong 服务端响应
type ServerPong struct {
}

// Login 登录
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
}

// LoginRes 登录结果
type LoginRes struct {
	ClientId string `json:"client_id"`
	Error    string `json:"error"`
}

type RegisterTunnel struct {
	Tunnels []tunnel.Tunnel
}

type RegisterTunnelRes struct {
	Tunnels []tunnel.Tunnel
	Error string `json:error`
}


