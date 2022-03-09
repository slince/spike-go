package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"net"
	"strconv"
)

type UdpHandler struct {
	proxyConnPool *conn.Pool
	conn  *net.UDPConn
}

func NewUdpHandler(connPool *conn.Pool) *UdpHandler{
	return &UdpHandler{
		proxyConnPool: connPool,
	}
}

func (udp *UdpHandler) Listen(serverPort int) error {

	var address, err = net.ResolveUDPAddr("udp", net.JoinHostPort("0.0.0.0", strconv.Itoa(serverPort)))
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", address)
	udp.conn = udpConn
	if err != nil {
		return err
	}
	var proxyConn = udp.proxyConnPool.Get()
	var bridge = cmd.NewBridge(proxyConn)

	go func() {
		buf := make([]byte, 1024)
		for {
			read, remoteAddr, _ := udpConn.ReadFromUDP(buf)
			if read == 0 {
				break
			}
			var udpPackage = &cmd.UdpPackage{
				Body: buf[0:read],
				RemoteAddr: remoteAddr,
			}
			err = bridge.Write(udpPackage)
			if err != nil {
				break
			}
		}
		udpConn.Close()
		proxyConn.Close()
	}()

	go func() {
		Handle:
		for {
			var command, err = bridge.Read()
			if err != nil {
				break
			}
			switch command := command.(type) {
			case *cmd.UdpPackage:
				_, err = udpConn.WriteToUDP(command.Body, command.RemoteAddr)
				if err != nil {
					break Handle
				}
			default:
				break Handle
			}
		}
		udpConn.Close()
		proxyConn.Close()
	}()
	return nil
}

func (udp *UdpHandler) AddProxyConn(proxyConn net.Conn) {
	udp.proxyConnPool.Put(proxyConn)
}

func (udp *UdpHandler) Close() {
	_ = udp.conn.Close()
}