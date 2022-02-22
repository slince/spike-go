package client

import (
	"fmt"
	"github.com/slince/spike/pkg/msg"
	"net"
	"strconv"
)

type Client struct {
	Id string
	Host string
	Port int
	Username string
	Password string
	Conn net.Conn
	Version string
}

func (cli *Client) Start() error{
	var address = cli.Host + ":" + strconv.Itoa(cli.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	cli.Conn = conn
	cli.login()
	err = cli.handleConn()
	return err
}

func (cli *Client) sendMsg(m interface{}) error{
	err := msg.WriteMsg(cli.Conn, m)
	return err
}

func (cli *Client) login(){
	cli.sendMsg(&msg.Login{
		Username: cli.Username,
		Password: cli.Password,
		Version: cli.Version,
	})
}

func (cli *Client) handleConn() error{
	for {
		rawMsg, err := msg.ReadMsg(cli.Conn)
		if err != nil {
			return err
		}
		switch m := rawMsg.(type) {
		case *msg.ServerPong:
		case *msg.LoginRes:
			if len(m.ClientId) > 0 {
				cli.Id = m.ClientId
			} else {
				return fmt.Errorf("login failed: %s", m.Error)
			}
		}
	}
}

func NewClient(host string, port int, username string, password string) *Client{
	return &Client{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Version:  "0.0.1",
	}
}