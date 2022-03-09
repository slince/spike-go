package proxy

import (
	"github.com/slince/spike/client"
	"github.com/slince/spike/pkg/conn"
	"github.com/slince/spike/pkg/log"
	"github.com/slince/spike/pkg/transfer"
	"net"
	"time"
)

type TcpHandler struct {
	logger *log.Logger
	cli *client.Client
	localAddress string
	proxyConn net.Conn
	bridge *transfer.Bridge
}

func NewTcpHandler() *TcpHandler{
	return &TcpHandler{

	}
}

func (tcp *TcpHandler) start(proxyConn net.Conn) {
	localConn, err := tcp.newLocalConn()
	if err != nil {
		return
	}
	conn.Combine(localConn, proxyConn)
	tcp.logger.Info("The worker is closed")
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

