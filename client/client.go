package client

import (
	"errors"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Id string
	Host string
	Port int
	User auth.GenericUser
	Conn net.Conn
	Bridge *transfer.Bridge
	Version string
	LastActiveAt time.Time
	logger *log.Logger
	handler *Handler
}

func NewClient(config Configuration) (*Client, error){
	var logger, err = log.NewLogger(config.Log)
	if err != nil {
		return nil, err
	}
	var cli = &Client{
		Host:     config.Host,
		Port:     config.Port,
		User: config.Auth,
		Version:  "0.0.1",
		LastActiveAt: time.Now(),
		logger: logger,
	}
	cli.handler = &Handler{
		client: cli,
		config: config,
	}
	return cli, err
}

func (cli *Client) Start() (err error){
	cli.Conn, err = cli.newConn()
	if err != nil {
		return
	}
	cli.Bridge = transfer.NewBridge(ft, cli.Conn, cli.Conn)
	err = cli.login()
	if err != nil {
		return
	}
	err = cli.handleConn()
	return
}

func (cli *Client) newConn() (net.Conn, error){
	var address = cli.Host + ":" + strconv.Itoa(cli.Port)
	conn, err := net.DialTimeout("tcp", address, 5)
	cli.logger.Info("Connected to the server")
	return conn, err
}

func (cli *Client) sendCommand(command transfer.Command) error{
	err := cli.Bridge.Write(command)
	cli.LastActiveAt = time.Now()
	return err
}

func (cli *Client) login() error {
	return cli.sendCommand(&cmd.Login{
		Username: cli.User.Username,
		Password: cli.User.Password,
		Version: cli.Version,
	})
}

func (cli *Client) handleConn() error{
	for {
		command, err := cli.Bridge.Read()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				err = errors.New("the connection is expired")
			}
			cli.logger.Warn("Failed to read command: ", err)
			return err
		}
		cli.logger.Trace("Receive a command:", command)
		switch command := command.(type) {
		case *cmd.ServerPong:
		case *cmd.LoginRes:
			if len(command.ClientId) > 0 {
				cli.Id = command.ClientId
				cli.logger.Info("The client is connected to the server, client id:", cli.Id)
				err = cli.handler.registerTunnels()
				if err != nil {
					return err
				}
			} else {
				cli.logger.Error("Failed to logged to the server: ", err)
				return err
			}
		case *cmd.RequestProxy:
			err = cli.handler.registerProxy(command)
			if err != nil {
				return err
			}
		}
	}
}