package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
	"sync"
)

type UdpHandler struct {
	logger *log.Logger
	proxyConnPool *conn.Pool
	tun tunnel.Tunnel

	listener  *net.UDPConn
	listenAddress string
	//stop chan bool
}

func NewUdpHandler(logger *log.Logger, connPool *conn.Pool, tun tunnel.Tunnel) *UdpHandler{
	return &UdpHandler{
		logger: logger,
		proxyConnPool: connPool,
		tun: tun,
		listenAddress: net.JoinHostPort("0.0.0.0", strconv.Itoa(tun.ServerPort)),
		//stop: make(chan bool, 1),
	}
}

func (udp *UdpHandler) Listen() (chan bool, error) {
	var address, err = net.ResolveUDPAddr("udp", udp.listenAddress)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenUDP("udp", address)
	if err != nil {
		return nil, err
	}
	udp.listener = listener

	var stop = make(chan bool, 1)

	var readUdp = func(proxyConn net.Conn, bridge *transfer.Bridge, wait *sync.WaitGroup) {
		defer listener.Close()
		defer proxyConn.Close()
		defer wait.Done()

		buf := make([]byte, 1024)
		for {
			read, remoteAddr, _ := listener.ReadFromUDP(buf)
			if read == 0 {
				stop <- true
				break
			}
			var udpPackage = &cmd.UdpPackage{
				Body: buf[0:read],
				RemoteAddr: remoteAddr,
			}
			err = bridge.Write(udpPackage)
			if err != nil {
				udp.logger.Error("Failed to write udp package to proxy conn:", err)
				break
			}
		}
	}

	var readProxy = func(proxyConn net.Conn, bridge *transfer.Bridge, wait *sync.WaitGroup) {
		defer proxyConn.Close()
		defer wait.Done()

		Handle:
		for {
			var command, err = bridge.Read()
			if err != nil {
				udp.logger.Error("Failed to read udp package from proxy conn, error: ", err)
				break
			}
			switch command := command.(type) {
			case *cmd.UdpPackage:
				_, err = listener.WriteToUDP(command.Body, command.RemoteAddr)
				if err != nil {
					udp.logger.Error("Failed to send udp package to pub conn ", err)
				}
			default:
				break Handle
			}
		}
	}

	var listenStop = make(chan bool, 1)
	var handleUdp = func() {
		for {
			select {
			case <- stop:
				listenStop <- true
				return
			default:
				var proxyConn, err = udp.proxyConnPool.Get()
				if err != nil {
					udp.logger.Error("The worker is closed, failed to get proxy conn from client, error: ", err)
					return
				}
				var bridge = cmd.NewBridge(proxyConn)
				var wait sync.WaitGroup
				wait.Add(2)
				go readUdp(proxyConn, bridge, &wait)
				go readProxy(proxyConn, bridge, &wait)
				wait.Wait()
			}
		}
	}
	go handleUdp()
	return listenStop, nil
}

func (udp *UdpHandler) AddProxyConn(proxyConn net.Conn) {
	udp.proxyConnPool.Put(proxyConn)
}

func (udp *UdpHandler) Close() {
	_ = udp.listener.Close()
}