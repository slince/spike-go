package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"strconv"
)

type UdpHandler struct {
	logger *log.Logger
	proxyConnPool *conn.Pool
	listener  *net.UDPConn
}

func NewUdpHandler(logger *log.Logger, connPool *conn.Pool) *UdpHandler{
	return &UdpHandler{
		logger: logger,
		proxyConnPool: connPool,
	}
}

func (udp *UdpHandler) Listen(serverPort int) error {

	var address, err = net.ResolveUDPAddr("udp", net.JoinHostPort("0.0.0.0", strconv.Itoa(serverPort)))
	if err != nil {
		return err
	}
	listener, err := net.ListenUDP("udp", address)
	udp.listener = listener
	if err != nil {
		return err
	}

	go func() {
		var proxyConn, err  = udp.proxyConnPool.Get()
		if err != nil {
			udp.logger.Error("The worker is closed, failed to get proxy conn from client, error: ", err)
			return
		}
		var bridge = cmd.NewBridge(proxyConn)

		go func() {
			buf := make([]byte, 1024)
			for {
				read, remoteAddr, _ := listener.ReadFromUDP(buf)
				if read == 0 {
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
			listener.Close()
			proxyConn.Close()
		}()

		go func() {
			var err error
			Handle:
			for {
				var command transfer.Command
				command, err = bridge.Read()
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

			listener.Close()
			proxyConn.Close()
		}()
	}()

	return nil
}

func (udp *UdpHandler) AddProxyConn(proxyConn net.Conn) {
	udp.proxyConnPool.Put(proxyConn)
}

func (udp *UdpHandler) Close() {
	_ = udp.listener.Close()
}