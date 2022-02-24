package client

import (
	"fmt"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/tunnel"
)

type Handler struct {
	client *Client
	config Configuration
}

func (h *Handler) registerTunnels() error{
	return h.client.sendCommand(&cmd.RegisterTunnel{
		ClientId: h.client.Id,
		Tunnels: h.config.Tunnels,
	})
}

func (h *Handler) registerProxy(command *cmd.RequestProxy) error{
	var tun, ok = getTunnel(h.config.Tunnels, command.ServerPort)
	if !ok {
		return fmt.Errorf("cannot find tunnel config for server port: %d", command.ServerPort)
	}
	var worker = newWorker(h.client, tun)
	go worker.start()
	return nil
}


func getTunnel(tunnels []tunnel.Tunnel, serverPort uint16) (tunnel.Tunnel, bool){
	var target tunnel.Tunnel
	var ok bool
	for _, tun := range tunnels {
		if tun.ServerPort == serverPort {
			target = tun
			ok = true
			break
		}
	}
	return target, ok
}