package proxy

import (
	"github.com/slince/spike/pkg/cmd"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"time"
)

type TcpHandler struct {
	logger *log.Logger
	localAddress string
	proxyConn net.Conn
	bridge *transfer.Bridge
}

func NewTcpHandler(logger *log.Logger, localAddress string, proxyConn net.Conn) *TcpHandler{
	return &TcpHandler{
		logger: logger,
		localAddress: localAddress,
		proxyConn: proxyConn,
		bridge: cmd.NewBridge(proxyConn),
	}
}

func (tcp *TcpHandler) Start() error {
	localConn, err := tcp.newLocalConn()
	if err != nil {
		return err
	}
	conn.Combine(localConn, tcp.proxyConn)
	return nil
}


func (tcp *TcpHandler) newLocalConn() (net.Conn, error){
	var con, err = net.DialTimeout("tcp", tcp.localAddress, 5 * time.Second)
	if err != nil {
		tcp.logger.Error("Failed to connect local service: ", err)
	} else {
		tcp.logger.Info("Connected to the local service: ", tcp.localAddress)
	}
	return con,err
}

