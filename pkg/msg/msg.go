package msg

import (
	"github.com/slince/spike/pkg/transfer"
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

func (s *ServerPong) GetType() transfer.MsgType {
	return TypePong
}

// Login 登录
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
}

func (s *Login) GetType() transfer.MsgType {
	return TypeLogin
}

// LoginRes 登录结果
type LoginRes struct {
	ClientId string `json:"client_id"`
	Error    string `json:"error"`
}

func (s *LoginRes) GetType() transfer.MsgType {
	return TypeLoginRes
}
