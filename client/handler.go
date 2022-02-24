package client

import (
	"fmt"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/tunnel"
)

func (cli *Client) registerTunnels() error{
	return cli.sendCommand(&cmd.RegisterTunnel{
		ClientId: cli.id,
		Tunnels: cli.tunnels,
	})
}

func (cli *Client) registerProxy(command *cmd.RequestProxy) error{
	var tun, ok = getTunnel(cli.tunnels, command.ServerPort)
	if !ok {
		return fmt.Errorf("cannot find tunnel config for server port: %d", command.ServerPort)
	}
	var worker = newWorker(cli, tun)
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