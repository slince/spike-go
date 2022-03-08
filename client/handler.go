package client

import (
	"errors"
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


func (cli *Client) handleRegisterTunnelRes(command *cmd.RegisterTunnelRes) error{
	var errNum = 0
	for _, result := range command.Results {
		if len(result.Error) > 0 {
			cli.logger.Warn(result.Error)
			errNum ++
		}
	}
	if errNum > 0 && errNum == len(command.Results) {
		return errors.New("all tunnels register failed")
	}
	return nil
}

func (cli *Client) handleViewProxyResp(command *cmd.ViewProxyResp){
	cli.proxiesChan <- command.Items
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


func getTunnel(tunnels []tunnel.Tunnel, serverPort int) (tunnel.Tunnel, bool){
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