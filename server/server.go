package server

import (
	"errors"
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

type Client struct {
	Id            string
	RemoteAddress string
	Conn          net.Conn
	Bridge        *transfer.Bridge
	ActiveAt      time.Time
	Tunnels       []tunnel.Tunnel
}

func NewClient(conn net.Conn, bridge *transfer.Bridge) *Client {
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
	Workers map[uint16]*Worker
	lock    sync.Mutex
	logger *log.Logger
}

func NewServer(cfg Configuration) (*Server, error) {
	var au = auth.NewSimpleAuth(cfg.Users)
	var logger, err = log.NewLogger(cfg.Log)
	if err != nil {
		return nil, err
	}
	var ser =  &Server{
		Host:    cfg.Host,
		Port:    cfg.Port,
		Clients: make(map[net.Conn]*Client, 0),
		Auth:    au,
		Workers: make(map[uint16]*Worker, 0),
		logger: logger,
	}
	return ser, nil
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
	ser.logger.Info("The server is running on " + address)
	for {
		conn, err := socket.Accept()
		if err != nil {
			ser.logger.Warn("Failed to accept connection: ", err)
			continue
		}
		ser.logger.Info("Accept a connection from : ", conn.RemoteAddr())
		go ser.handleConn(conn)
	}
}

func (ser *Server) handleConn(conn net.Conn) {
	var bridge = transfer.NewBridge(ft, conn, conn)

	ParseCommand:
		for {
			var command, err = bridge.Read()
			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					err = errors.New("the client connection is expired")
				}
				ser.logger.Warn("Failed to read command: ", err)
				if client, ok := ser.Clients[conn]; ok {
					ser.closeClient(client)
				}
				err = conn.Close()
				break
			}
			ser.logger.Trace("Receive a command:", command)
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
				ser.logger.Print("handle register ", stop, err)
				if stop { // stop listen the socket.
					break ParseCommand
				}
			}
			if err != nil {
				ser.logger.Warn("Handle command error: ", err)
				_ = conn.Close()
			}
		}
}

func (ser *Server) sendCommand(client *Client, command transfer.Command) error {
	client.ActiveAt = time.Now()
	return client.Bridge.Write(command)
}

func (ser *Server) closeClient(client *Client) {
	defer ser.lock.Unlock()
	ser.lock.Lock()
	for _, tun := range client.Tunnels {
		ser.closeTunnel(tun)
	}
	delete(ser.Clients, client.Conn)
}

func (ser *Server) closeTunnel(tun tunnel.Tunnel) {
	if worker, ok := ser.Workers[tun.ServerPort]; ok {
		_ = worker.Close()
		delete(ser.Workers, tun.ServerPort)
	}
}
