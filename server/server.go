package server

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
	"sync"
	"time"
)

var logger = log.NewLogger()

func init() {
	logger.EnableConsole()
}

type Client struct {
	Id            string
	RemoteAddress string
	Conn          net.Conn
	Bridge        *transfer.Bridge
	ActiveAt      time.Time
	Tunnels       []tunnel.Tunnel
}

func NewClient(conn net.Conn, bridge *transfer.Bridge) *Client {
	//var bridge = transfer.NewBridge(ft, conn, conn)

	return &Client{
		xid.New().String(),
		conn.RemoteAddr().String(),
		conn,
		bridge,
		time.Now(),
		make([]tunnel.Tunnel, 0),
	}
}

type Server struct {
	Host    string
	Port    int
	Clients map[net.Conn]*Client
	Auth    auth.Auth
	Workers map[tunnel.Tunnel]*Worker
	lock    sync.Mutex
}

func NewServer(host string, port int, au auth.Auth) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Clients: make(map[net.Conn]*Client, 0),
		Auth:    au,
		Workers: make(map[tunnel.Tunnel]*Worker, 0),
	}
}

func (ser *Server) GetClient(id string) (*Client, error) {
	for _, client := range ser.Clients {
		if client.Id == id {
			return client, nil
		}
	}
	return nil, fmt.Errorf("cannot find client with id %s", id)
}

func (ser *Server) Start() error {
	var address = ser.Host + ":" + strconv.Itoa(ser.Port)
	socket, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	logger.Info("The server is running on " + address)
	for {
		conn, err := socket.Accept()
		if err != nil {
			logger.Warn("Failed to accept connection: ", err)
			continue
		}
		logger.Info("Accept a connection from : ", conn.RemoteAddr())
		go ser.handleConn(conn)
	}
}

func (ser *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	var bridge = transfer.NewBridge(ft, conn, conn)
	for {
		command, err := bridge.Read()
		if err != nil {
			logger.Warn("Failed to read command: ", err)
			if client, ok := ser.Clients[conn]; ok {
				err = ser.closeClient(client)
				if err != nil {
					return
				}
			}
			err = conn.Close()
			if err != nil {
				panic(err)
			}
			return
		}
		logger.Trace("Receive a command:", command)
		switch command := command.(type) {
		case *cmd.ClientPing:
			err = ser.handlePing(command)
		case *cmd.Login:
			err = ser.handleLogin(command, conn, bridge)
		case *cmd.RegisterTunnel:
			err = ser.handleRegisterTun(command, conn, bridge)
		case *cmd.RegisterProxy:
			var stop = false
			stop, err = ser.handleRegisterProxy(command, conn, bridge)
			if stop { // stop listen the socket.
				break
			}
		}
	}
}

func (ser *Server) sendCommand(client *Client, command transfer.Command) error {
	client.ActiveAt = time.Now()
	return client.Bridge.Write(command)
}

func (ser *Server) handlePing(m *cmd.ClientPing) error {
	var client, err = ser.GetClient(m.ClientId)
	if err != nil {
		return err
	}
	client.ActiveAt = time.Now()
	return ser.sendCommand(client, cmd.ServerPong{})
}

func (ser *Server) handleLogin(command *cmd.Login, conn net.Conn, bridge *transfer.Bridge) (err error) {
	defer ser.lock.Unlock()
	if user := ser.Auth.Check(command); user != nil {
		client := NewClient(conn, bridge)
		ser.lock.Lock()
		ser.Clients[conn] = client
		err = ser.sendCommand(client, cmd.LoginRes{ClientId: client.Id})
	} else {
		err = bridge.Write(cmd.LoginRes{ClientId: "", Error: "cannot verify your identify"})
	}
	return
}

func (ser *Server) handleRegisterTun(command *cmd.RegisterTunnel, conn net.Conn, bridge *transfer.Bridge) error {
	defer ser.lock.Unlock()
	var client, err = ser.GetClient(command.ClientId)
	if err != nil {
		return err
	}
	ser.lock.Lock()
	for _, tun := range command.Tunnels {
		ser.Workers[tun] = newWorker(tun, conn, bridge)
		logger.Info("Starting the worker for tunnel ", tun.Id)
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
		worker.AddProxyConn(conn)
		return true, nil
	}
	return false, fmt.Errorf("cannot find worker for the tunnel %s", command.Tunnel.Id)
}

func (ser *Server) closeClient(client *Client) error {
	defer ser.lock.Unlock()
	for _, tun := range client.Tunnels {
		err := ser.closeTunnel(tun)
		if err != nil {
			return err
		}
	}
	ser.lock.Lock()
	delete(ser.Clients, client.Conn)
	return nil
}

func (ser *Server) closeTunnel(tun tunnel.Tunnel) error {
	defer ser.lock.Unlock()
	ser.lock.Lock()
	if worker, ok := ser.Workers[tun]; ok {
		err := worker.Close()
		if err != nil {
			return err
		}
		delete(ser.Workers, tun)
	}
	return nil
}
