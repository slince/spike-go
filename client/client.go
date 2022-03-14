package client

import (
	"errors"
	"fmt"
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	proxiesChan chan []cmd.ProxyItem
	onLogin func(command *cmd.LoginRes) error
}

func NewClient(config Configuration) (*Client, error){
	var logger, err = log.NewLogger(config.Log)
	if err != nil {
		return nil, err
	}
	var cli = &Client{
		host:     config.Host,
		port:     config.Port,
		user:     config.User,
		version:  "0.0.1",
		activeAt: time.Now(),
		logger:   logger,
		tunnels:  config.Tunnels,
		proxiesChan: make(chan []cmd.ProxyItem, 2),
		onLogin: func(command *cmd.LoginRes) error {
			return nil
		},
	}
	return cli, err
}

func (cli *Client) Start() (err error){
	cli.conn, err = cli.NewConn()
	if err != nil {
		return
	}
	cli.bridge = cmd.NewBridge(cli.conn)
	return
}

func (cli *Client) StartWithPrevSession() (err error){
	var prevSession string
	prevSession, err = attemptGetPrevSessionId()
	if err != nil {
		return
	}

	err = cli.Start()
	if err != nil {
		return
	}

	cli.id = prevSession
	return
}

func (cli *Client) Listen() (err error){
	err = cli.Start()
	if err != nil {
		return
	}
	cli.onLogin = cli.onLoginCall
	err = cli.login()
	if err != nil {
		return
	}
	go cli.graceExit()
	err = cli.handleConn()
	if err != nil {
		cli.logger.Error("Error: ", err)
	}
	return
}

func (cli *Client) onLoginCall(command *cmd.LoginRes) error{
	var err error
	if len(command.ClientId) > 0 {
		cli.id = command.ClientId
		cli.logger.Info("Logged in to the server, client id:", cli.id)
		var err2 = saveSessionId(cli.id)
		if err2 != nil {
			cli.logger.Warn("Fail to dump client id to the session file")
		}
		go cli.autoPing() // heartbeat
		err = cli.registerTunnels()
	} else {
		err= fmt.Errorf("failed to log in to the server: %s", command.Error)
	}
	return err
}

func (cli *Client) autoPing(){
	var timer = time.NewTicker(10 * time.Second)
	for range timer.C{
		_ = cli.sendCommand(&cmd.ClientPing{
			ClientId: cli.id,
		})
	}
}

func (cli *Client) NewConn() (net.Conn, error){
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
	defer func() {
		_ = removeSessionFile()
	}()
	for {
		var command transfer.Command
		command, err = cli.bridge.Read()
		if err != nil {
			if _, ok := err.(*net.OpError); ok || err == io.EOF{
				err = errors.New("the connection is expired")
			}
			cli.logger.Warn("Failed to read command from server: ", err)
			return
		}
		cli.logger.Trace("Receive a command:", command)
		switch command := command.(type) {
		case *cmd.ServerPong:
			cli.activeAt = time.Now()
		case *cmd.LoginRes:
			err = cli.onLogin(command)
		case *cmd.RegisterTunnelRes:
			err = cli.handleRegisterTunnelRes(command)
		case *cmd.RequestProxy:
			err = cli.registerProxy(command)
		case *cmd.ViewProxyResp:
			cli.handleViewProxyResp(command)
		}
		if err!= nil {
			_ = cli.conn.Close()
			return
		}
	}
}


func (cli *Client) GetProxies() ([]cmd.ProxyItem, error){
	var err = cli.StartWithPrevSession()
	if err != nil {
		cli.logger.Warn("Failed to login the server using prev session id, connect again")
		cli.onLogin = func(command *cmd.LoginRes) error {
			cli.id = command.ClientId
			return cli.sendCommand(&cmd.ViewProxy{
				ClientId: cli.id,
			})
		}
		err = cli.Start()
		if err != nil {
			return nil, err
		}
		err = cli.login()
		if err != nil {
			return nil, err
		}
	} else {
		err = cli.sendCommand(&cmd.ViewProxy{
			ClientId: cli.id,
		})
		if err != nil {
			return nil, err
		}
	}
	go cli.handleConn()
	var timer = time.After(5 * time.Second)
	for {
		select {
		case <-timer:
			return nil, errors.New("timeout to get proxies")
		case tunnels := <-cli.proxiesChan:
			return tunnels, nil
		}
	}
}

func (cli *Client) graceExit(){
	var exit = make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<- exit
	_ = removeSessionFile()
	os.Exit(0)
}

func saveSessionId(clientId string) error{
	var file, err = getSessionFile()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(clientId), fs.ModePerm)
}

func removeSessionFile() error{
	var file, err = getSessionFile()
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return err
	}
	return os.Remove(file)
}

func getSessionFile() (string, error){
	var dir, err = os.Getwd()
	if err != nil {
		return "", err
	}
	return dir + "/spike.session", nil
}

func attemptGetPrevSessionId() (string, error){
	var file, err = getSessionFile()
	if err != nil {
		return "", err
	}
	read, err := ioutil.ReadFile(file)
	return string(read), err
}
