package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"sync"
	"time"
)

type UdpHandler struct {
	logger *log.Logger
	localAddress string
	lock sync.Mutex
	localConnMap map[*net.UDPAddr]UdpConnUnit
	proxyConn net.Conn
	bridge *transfer.Bridge
	messages chan *cmd.UdpPackage
}

type UdpConnUnit struct {
	udpConn *net.UDPConn
	activeAt time.Time
}

func NewUdpHandler(logger *log.Logger, localAddress string, proxyConn net.Conn) *UdpHandler{
	return &UdpHandler{
		logger: logger,
		localAddress: localAddress,
		localConnMap: make(map[*net.UDPAddr]UdpConnUnit),
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
	go udp.checkAlive()

	for msg := range udp.messages {
		go udp.handleMessage(msg)
	}
	return nil
}

func (udp *UdpHandler) checkAlive(){
	var timer = time.NewTicker(10 * time.Second)
	defer timer.Stop()
	for range timer.C {
		udp.lock.Lock()
		var aliveTime = time.Now().Add(-20 * time.Second)
		for addr, unit := range udp.localConnMap {
			if unit.activeAt.After(aliveTime) {
				continue
			}
			_ = unit.udpConn.Close()
			delete(udp.localConnMap, addr)
		}
		udp.lock.Unlock()
	}
}

func (udp *UdpHandler) getLocalConn(remoteAddr *net.UDPAddr) (*net.UDPConn, error){
	udp.lock.Lock()
	defer udp.lock.Unlock()
	var err error
	var localConn *net.UDPConn
	var localConnUnit, ok = udp.localConnMap[remoteAddr]
	if !ok {
		localConn, err = udp.newLocalConn()
		if err != nil {
			return nil, err
		}
		localConnUnit = UdpConnUnit{
			udpConn: localConn,
			activeAt: time.Now(),
		}
		udp.localConnMap[remoteAddr] = localConnUnit
		go udp.joinLocalToProxy(localConn, remoteAddr)
	} else {
		localConn = localConnUnit.udpConn
		localConnUnit.activeAt = time.Now()
	}
	return localConn, nil
}

func (udp *UdpHandler) handleMessage(msg *cmd.UdpPackage) error{
	var localConn, err = udp.getLocalConn(msg.RemoteAddr)
	if err != nil {
		return err
	}
	_, err = localConn.Write(msg.Body)
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
		read, err := localConn.Read(buf)
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