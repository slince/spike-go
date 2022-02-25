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
	TypeRegisterTunnel
	TypeRegisterTunnelRes
	TypeRequestProxy
	TypeRegisterProxy
	TypeStartProxy
)

type ClientPing struct {
	transfer.BaseCommand
	ClientId string `json:"client_id"`
}

type ServerPong struct {
	transfer.BaseCommand
}

type Login struct {
	transfer.BaseCommand
	Username string
	Password string
	Version  string
}

type LoginRes struct {
	transfer.BaseCommand
	ClientId string
	Error    string
}

type RegisterTunnel struct {
	transfer.BaseCommand
	ClientId string
	Tunnels  []tunnel.Tunnel
}

type RegisterTunnelRes struct {
	transfer.BaseCommand
	Results []TunnelResult
}

func (r RegisterTunnelRes) AddResult(result TunnelResult){
	r.Results = append(r.Results, result)
}

type TunnelResult struct {
	Tun tunnel.Tunnel
	Error string
}

type RequestProxy struct {
	transfer.BaseCommand
	ServerPort uint16
}

type RegisterProxy struct {
	transfer.BaseCommand
	ClientId string
	Tunnel   tunnel.Tunnel
}

type StartProxy struct {
	transfer.BaseCommand
	Tunnel tunnel.Tunnel
}
