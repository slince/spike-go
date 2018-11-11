package server

import (
	"errors"
	"fmt"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type ChunkServer interface {
	//启动监听服务
	Run() error
	//获取对应的tunnel
	GetTunnel() tunnel.Tunnel
	// 设置代理连接
	SetProxyConnection(pubConnId string, conn net.Conn) error
}

// 监听公网接口
type TcpChunkServer struct {
	//对应的隧道
	Tunnel *tunnel.TcpTunnel
	//服务的客户端
	Client *Client
	//服务调度程序
	Server *Server
	//公共请求
	pubConnCollection map[string]*PublicConn
}

// 启动服务
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
		chunkServer.pubConnCollection[publicConnection.Id] = publicConnection
		//处理请求
		go chunkServer.handleConnection(publicConnection)
	}
}

func (chunkServer *TcpChunkServer) handleConnection(pubConn *PublicConn) {
	//1.收到公网请求，请求客户端代理
	msg := protol.Protocol{
		Action: "request_proxy",
		Headers: map[string]string{
			"tunnel-id": chunkServer.Tunnel.Id,
		},
	}
	chunkServer.Server.SendMessageToClient(chunkServer.Client, &msg)

	// 2. 挂起当前公网请求
	//var proxyConn net.Conn
	proxyConn := <- pubConn.ProxyConnChan //从通道读取代理请求
	defer close(pubConn.ProxyConnChan)

	// 3. 管道请求
	pubConn.Pipe(proxyConn)
	delete(chunkServer.pubConnCollection, pubConn.Id)
}

// 设置代理链接
func (chunkServer *TcpChunkServer) SetProxyConnection(pubConnId string, conn net.Conn) error{

	if pubConn, ok := chunkServer.pubConnCollection[pubConnId];ok {
		pubConn.ProxyConnChan <- conn
		return nil
	}
	return fmt.Errorf(`the public connection id "%s" is missing`, pubConnId)
}

func (chunkServer *TcpChunkServer) GetTunnel() tunnel.Tunnel{
	return chunkServer.Tunnel
}

// http chunk server
type HttpChunkServer struct {
	TcpChunkServer
}