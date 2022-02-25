package server

import (
	"fmt"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"time"
)

func (ser *Server) handlePing(m *cmd.ClientPing) error {
	var client, err = ser.GetClient(m.ClientId)
	if err != nil {
		return err
	}
	client.ActiveAt = time.Now()
	return ser.sendCommand(client, &cmd.ServerPong{})
}

func (ser *Server) handleLogin(command *cmd.Login, conn net.Conn, bridge *transfer.Bridge) (err error) {
	if user := ser.Auth.Check(command); user != nil {
		client := NewClient(conn, bridge)
		defer ser.lock.Unlock()
		ser.lock.Lock()
		ser.Clients[conn] = client
		err = ser.sendCommand(client, &cmd.LoginRes{ClientId: client.Id})
	} else {
		err = bridge.Write(&cmd.LoginRes{ClientId: "", Error: "cannot verify your identify"})
	}
	return
}

func (ser *Server) handleRegisterTun(command *cmd.RegisterTunnel, conn net.Conn, bridge *transfer.Bridge) error {
	var client, err = ser.GetClient(command.ClientId)
	if err != nil {
		return err
	}
	ser.lock.Lock()
	var result = &cmd.RegisterTunnelRes{}
	for _, tun := range command.Tunnels {
		var tunResult = cmd.TunnelResult{Tunnel: tun}
		if _, exists := ser.Workers[tun.ServerPort];exists {
			tunResult.Error = fmt.Sprintf("the tunnel for port %d is exists", tun.ServerPort)
		} else {
			ser.Workers[tun.ServerPort] = newWorker(ser, tun, conn, bridge)
			ser.logger.Info("Starting the worker for tunnel ", tun.ServerPort)
			err = ser.Workers[tun.ServerPort].Start()
			if err != nil {
				tunResult.Error = fmt.Sprint("Failed to start the worker", err.Error())
			} else {
				client.Tunnels = append(client.Tunnels, tun)
			}
		}
		if len(tunResult.Error) > 0 {
			ser.logger.Warn(tunResult.Error)
		}
		result.AddResult(tunResult)
	}
	ser.lock.Unlock()
	return ser.sendCommand(client, result)
}

func (ser *Server) handleRegisterProxy(command *cmd.RegisterProxy, conn net.Conn, bridge *transfer.Bridge) (bool, error) {
	var _, err = ser.GetClient(command.ClientId)
	if err != nil {
		return false, err
	}
	if worker, ok := ser.Workers[command.Tunnel.ServerPort]; ok {
		worker.addProxyConn(conn)
		return true, nil
	}
	return false, fmt.Errorf("cannot find worker for the tunnel port %d", command.Tunnel.ServerPort)
}
