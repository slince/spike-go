package server

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Id            string
	RemoteAddress string
	Conn          net.Conn
	Bridge *transfer.Bridge
	ActiveAt      time.Time
}

func NewClient(conn net.Conn, bridge *transfer.Bridge) *Client {
	//var bridge = transfer.NewBridge(ft, conn, conn)

	return &Client{
		xid.New().String(),
		conn.RemoteAddr().String(),
		conn,
		bridge,
		time.Now(),
	}
}

type Server struct {
	Host    string
	Port    int
	Clients []*Client
	Auth    auth.Auth
}

var logger = log.NewLogger()

func init()  {
	logger.EnableConsole()
}

func (ser *Server) GetClient(id string) *Client {
	for _, client := range ser.Clients {
		if client.Id == id {
			return client
		}
	}
	return nil
}

func (ser *Server) Start() error {
	var address = ser.Host + ":" + strconv.Itoa(ser.Port)
	socket, err := net.Listen("tcp", address)
	if err != nil {
		ser.handleError()
		return err
	}
	logger.Info("The server is running on " + address)
	for {
		conn, err := socket.Accept()
		if err != nil {
			logger.Warn("Failed to accept connection: ", err)
			continue
		}
		go ser.handleConn(conn)
	}
}

func (ser *Server) handleConn(conn net.Conn) {
	var bridge = transfer.NewBridge(ft, conn, conn)
	for {
		command, err := bridge.Read()
		if err != nil {
			logger.Warn("Failed to read command: ", err)
			conn.Close()
			return
		}
		switch command := command.(type) {
		case *cmd.ClientPing:
			err = ser.handlePing(command)
		case *cmd.Login:
			err = ser.handleLogin(command, conn, bridge)
		}
	}
}

func (ser *Server) sendCommand(client *Client, command transfer.Command) error {
	client.ActiveAt = time.Now()
	err := client.Bridge.Write(command)
	if err != nil {
		return nil
	}
	return err
}

func (ser *Server) handlePing(m *cmd.ClientPing) error {
	var client = ser.GetClient(m.ClientId)
	if client == nil {
		return fmt.Errorf("cannot find client with id %s", m.ClientId)
	}
	client.ActiveAt = time.Now()
	err := ser.sendCommand(client, cmd.ServerPong{})
	if err != nil {
		return err
	}
	return nil
}

func (ser *Server) handleLogin(command *cmd.Login, conn net.Conn, bridge *transfer.Bridge) (err error) {
	if user := ser.Auth.Check(command); user != nil {
		client := NewClient(conn, bridge)
		ser.Clients = append(ser.Clients, client)
		err = ser.sendCommand(client, cmd.LoginRes{ClientId: client.Id})
	} else {
		err = bridge.Write(cmd.LoginRes{ClientId: "", Error: "cannot verify your identify"})
	}
	return
}

func (ser *Server) handleError() {

}

func (ser *Server) handleRegisterTun() {

}

func NewServer(host string, port int, au auth.Auth) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Clients: make([]*Client, 5),
		Auth:    au,
	}
}
