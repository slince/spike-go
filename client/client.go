package client

import (
	"errors"
	"fmt"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
	"time"
)

type Client struct {
	id   string
	host string
	port int
	user auth.GenericUser
	conn net.Conn
	bridge *transfer.Bridge
	version  string
	activeAt time.Time
	logger   *log.Logger
	tunnels []tunnel.Tunnel
}

func NewClient(config Configuration) (*Client, error){
	var logger, err = log.NewLogger(config.Log)
	if err != nil {
		return nil, err
	}
	var cli = &Client{
		host:     config.Host,
		port:     config.Port,
		user:     config.Auth,
		version:  "0.0.1",
		activeAt: time.Now(),
		logger:   logger,
		tunnels:  config.Tunnels,
	}
	return cli, err
}

func (cli *Client) Start() (err error){
	cli.conn, err = cli.newConn()
	if err != nil {
		return
	}
	cli.bridge = transfer.NewBridge(ft, cli.conn, cli.conn)
	err = cli.login()
	if err != nil {
		return
	}
	err = cli.handleConn()
	if err != nil {
		cli.logger.Error("Error: ", err)
	}
	return
}

func (cli *Client) newConn() (net.Conn, error){
	var address = cli.host + ":" + strconv.Itoa(cli.port)
	conn, err := net.DialTimeout("tcp", address, 5 * time.Second)
	if err == nil {
		cli.logger.Info("Connected to the server")
	} else {
		cli.logger.Error("Failed to connect the server: ", err)
	}
	return conn, err
}

func (cli *Client) sendCommand(command transfer.Command) error{
	err := cli.bridge.Write(command)
	cli.activeAt = time.Now()
	return err
}

func (cli *Client) login() error {
	return cli.sendCommand(&cmd.Login{
		Username: cli.user.Username,
		Password: cli.user.Password,
		Version: cli.version,
	})
}

func (cli *Client) handleConn() (err error){
	for {
		var command transfer.Command
		command, err = cli.bridge.Read()
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				err = errors.New("the connection is expired")
			}
			cli.logger.Warn("Failed to read command: ", err)
			return
		}
		cli.logger.Trace("Receive a command:", command)
		switch command := command.(type) {
		case *cmd.ServerPong:
			cli.activeAt = time.Now()
		case *cmd.LoginRes:
			if len(command.ClientId) > 0 {
				cli.id = command.ClientId
				cli.logger.Info("Logged in to the server, client id:", cli.id)
				err = cli.registerTunnels()
			} else {
				err= fmt.Errorf("failed to log in to the server: %s", command.Error)
			}
		case *cmd.RegisterTunnelRes:
			err = cli.handleRegisterTunnelRes(command)
		case *cmd.RequestProxy:
			err = cli.registerProxy(command)
		}
		if err!= nil {
			_ = cli.conn.Close()
			return
		}
	}
}