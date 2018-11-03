package chunk_server

import (
	"bufio"
	"errors"
	"github.com/slince/jinbox/protol"
	"github.com/slince/jinbox/server"
	"github.com/slince/jinbox/tunnel"
	"net"
	"strings"
)

type ChunkServer interface {
	//启动监听服务
	Run()
	//获取对应的tunnel
	GetTunnel() tunnel.Tunnel
	// 设置代理连接
	SetProxyConnection(publicConnectionId string, conn net.Conn)
}

// 监听公网接口
type TcpChunkServer struct {
	//对应的隧道
	Tunnel *tunnel.TcpTunnel
	//服务的客户端
	Client *server.Client
	//服务调度程序
	Server *server.Server
}

func (chunkServer *TcpChunkServer) Run() error {
	ln,err := net.Listen("tcp", ":" + chunkServer.Tunnel.ServerPort)
	if  err != nil {
		return errors.New("failed to create chunk server")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		publicConnection := NewPublicConn(conn)
		go chunkServer.handleConnection(publicConnection)
	}
}

func (chunkServer *TcpChunkServer) handleConnection(conn *PublicConn) {
	//1.收到公网请求，请求客户端代理
	msg := protol.Protocol{
		Action: "request_proxy",
		Headers: map[string]string{
			"tunnel-id": chunkServer.Tunnel.Id,
		},
	}
	chunkServer.Server.SendMessageToClient(chunkServer.Client, &msg)

	bufio.NewReader(conn.Conn).ReadBytes()


	// 2. 挂起当前公网请求
	var proxyConn net.Conn
	a <- conn.ProxyConnChan

	rconn, wconn := net.Pipe()

}

func (chunkServer *TcpChunkServer) GetTunnel() tunnel.Tunnel{
	return chunkServer.Tunnel
}

// http chunk server
type HttpChunkServer struct {
	TcpChunkServer
}