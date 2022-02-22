package client

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"strconv"
	"time"
)

var logger = log.NewLogger()

func init()  {
	logger.EnableConsole()
}

type Client struct {
	Id string
	Host string
	Port int
	Username string
	Password string
	Conn net.Conn
	Bridge *transfer.Bridge
	Version string
	LastActiveAt time.Time
}

func (cli *Client) Start() (err error){
	var address = cli.Host + ":" + strconv.Itoa(cli.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	cli.Conn = conn
	cli.Bridge = transfer.NewBridge(ft, conn, conn)
	err = cli.login()
	if err != nil {
		return
	}
	err = cli.handleConn()
	return
}

func (cli *Client) sendCommand(command transfer.Command) error{
	err := cli.Bridge.Write(command)
	cli.LastActiveAt = time.Now()
	return err
}

func (cli *Client) login() error {
	return cli.sendCommand(cmd.Login{
		Username: cli.Username,
		Password: cli.Password,
		Version: cli.Version,
	})
}

func (cli *Client) handleConn() error{
	for {
		command, err := cli.Bridge.Read()
		if err != nil {
			return err
		}
		logger.Trace("Receive a command:", command)
		switch command := command.(type) {
		case *cmd.ServerPong:
		case *cmd.LoginRes:
			if len(command.ClientId) > 0 {
				cli.Id = command.ClientId
				logger.Info("The client is connected to the server, client id:", cli.Id)
			} else {
				logger.Error("Failed to logged to the server: ", err)
				return err
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