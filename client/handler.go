package client

import "github.com/slince/spike/pkg/cmd"

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

func (h *Handler) registerProxy() error{
	return h.client.sendCommand(&cmd.RegisterTunnel{
		ClientId: h.client.Id,
		Tunnels: h.config.Tunnels,
	})
}