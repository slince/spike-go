package server

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/msg"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Id            string
	RemoteAddress string
	Conn          net.Conn
	ActiveAt      time.Time
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		xid.New().String(),
		conn.RemoteAddr().String(),
		conn,
		time.Now(),
	}
}

type Server struct {
	Host    string
	Port    int
	Clients []*Client
	Auth    auth.Auth
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
	for {
		conn, err := socket.Accept()
		if err != nil {

		}
		go ser.handleConn(conn)
	}
}

func (ser *Server) handleConn(conn net.Conn) {
	for {
		rawMsg, err := msg.ReadMsg(conn)
		fmt.Println(rawMsg, err, rawMsg.(*msg.Login))
		if err != nil {
			return
		}
		switch m := rawMsg.(type) {
		case *msg.ClientPing:
			err = ser.handlePing(m)
		case *msg.Login:
			err = ser.handleLogin(m, conn)
		}
	}
}

func (ser *Server) sendMsg(client *Client, m interface{}) {
	client.ActiveAt = time.Now()
	msg.WriteMsg(client.Conn, m)
}

func (ser *Server) handlePing(m *msg.ClientPing) error {
	var client = ser.GetClient(m.ClientId)
	if client == nil {
		return fmt.Errorf("cannot find client with id %s", m.ClientId)
	}
	client.ActiveAt = time.Now()
	return nil
}

func (ser *Server) handleLogin(cmd *msg.Login, conn net.Conn) error {
	if user := ser.Auth.Check(cmd); user != nil {
		client := NewClient(conn)
		ser.Clients = append(ser.Clients, client)
		ser.sendMsg(client, &msg.LoginRes{ClientId: client.Id})
		fmt.Println("Login ok")
	} else {
		msg.WriteMsg(conn, &msg.LoginRes{ClientId: "", Error: "cannot verify your identify"})
		fmt.Println("Login fail")
	}
	return nil
}

func (ser *Server) handleError() {

}

func (ser *Server) handleRegisterTun(cmd msg.RegisterTun) {
	for _, tun := range cmd.Tunnels {

	}
}

func NewServer(host string, port int, au auth.Auth) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Clients: make([]*Client, 5),
		Auth:    au,
	}
}
