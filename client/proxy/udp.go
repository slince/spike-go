package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"sync"
)

type UdpHandler struct {
	logger *log.Logger
	localAddress string
	lock sync.Mutex
	localConnMap map[*net.UDPAddr]*net.UDPConn
	proxyConn net.Conn
	bridge *transfer.Bridge
	messages chan *cmd.UdpPackage
}

func NewUdpHandler(logger *log.Logger, localAddress string, proxyConn net.Conn) *UdpHandler{
	return &UdpHandler{
		logger: logger,
		localAddress: localAddress,
		localConnMap: make(map[*net.UDPAddr]*net.UDPConn),
		proxyConn: proxyConn,
		bridge: cmd.NewBridge(proxyConn),
		messages: make(chan *cmd.UdpPackage, 100),
	}
}

func (udp *UdpHandler) Start() error{
	// Read Msg
	go func() {
		Handle:
		for {
			var command, err = udp.bridge.Read()
			if err != nil {
				break
			}
			switch command := command.(type) {
			case *cmd.UdpPackage:
				udp.messages <- command
			default:
				break Handle
			}
		}
		udp.proxyConn.Close()
	}()

	for msg := range udp.messages {
		go udp.handleMessage(msg)
	}
	return nil
}

func (udp *UdpHandler) handleMessage(msg *cmd.UdpPackage) error{
	udp.lock.Lock()
	localConn, ok := udp.localConnMap[msg.RemoteAddr]
	var err error
	if !ok {
		localConn, err = udp.newLocalConn()
		if err != nil {
			return err
		}
		udp.localConnMap[msg.RemoteAddr] = localConn
		go udp.joinLocalToProxy(localConn, msg.RemoteAddr)
	}
	udp.lock.Unlock()
	_, err = localConn.WriteToUDP(msg.Body, nil)
	if err != nil {
		udp.lock.Lock()
		delete(udp.localConnMap, msg.RemoteAddr)
		localConn.Close()
		udp.lock.Unlock()
	}
	return err
}

func (udp *UdpHandler) joinLocalToProxy(localConn *net.UDPConn, remoteAddr *net.UDPAddr){
	var buf = make([]byte, 1024)
	for {
		read, _, err := localConn.ReadFromUDP(buf)
		if read == 0 {
			break
		}
		var udpPackage = &cmd.UdpPackage{
			Body: buf[0:read],
			RemoteAddr: remoteAddr,
		}
		err = udp.bridge.Write(udpPackage)
		if err != nil {
			break
		}
	}
	_ = udp.proxyConn.Close()
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