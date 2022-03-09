package proxy

import (
	"github.com/slince/spike/client"
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"net"
)

type UdpHandler struct {
	logger *log.Logger
	cli *client.Client
	localAddress string
	conn *net.UDPConn
}

func (udp *UdpHandler) Start(proxyConn net.Conn) error{
	var conn , err = udp.newLocalConn()
	if err != nil {
		return err
	}
	var bridge = cmd.NewBridge(proxyConn)
	go func() {
		var buf = make([]byte, 1024)
		for {
			read, remoteAddr, err := conn.ReadFromUDP(buf)
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
	}()
	return nil
}

func (udp *UdpHandler) newLocalConn() (*net.UDPConn, error){
	var address, err = net.ResolveUDPAddr("udp", udp.localAddress)
	if err != nil {
		return nil, err
	}
	con, err := net.DialUDP("udp", nil, address)
	if err != nil {
		udp.logger.Error("Failed to connect local service: ", err)
	} else {
		udp.logger.Info("Connected to the local service: ", udp.localAddress)
	}
	return con,err
}