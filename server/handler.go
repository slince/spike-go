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
	defer ser.lock.Unlock()
	ser.lock.Lock()
	for _, tun := range command.Tunnels {
		ser.Workers[tun] = newWorker(ser, tun, conn, bridge)
		ser.logger.Info("Starting the worker for tunnel ", tun.Id)
		err = ser.Workers[tun].Start()
		if err != nil {
			return err
		}
	}
	client.Tunnels = command.Tunnels
	return nil
}

func (ser *Server) handleRegisterProxy(command *cmd.RegisterProxy, conn net.Conn, bridge *transfer.Bridge) (bool, error) {
	var _, err = ser.GetClient(command.ClientId)
	if err != nil {
		return false, err
	}
	if worker, ok := ser.Workers[command.Tunnel]; ok {
		worker.addProxyConn(conn)
		return true, nil
	}
	return false, fmt.Errorf("cannot find worker for the tunnel %s", command.Tunnel.Id)
}
