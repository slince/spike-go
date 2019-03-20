package client

import (
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type Worker interface {
	// Start the worker
	Start() error
}

type TcpWorker struct {
	client *Client
	publicConnId string
	localConn net.Conn
	proxyConn *ProxyConn
	tunnel tunnel.Tunnel
}

func (worker *TcpWorker) Start() error{
	conn, err := worker.client.createConnector()
	if err != nil {
		return err
	}
	// 发送消息给控制服务
	msg := &protol.Protocol{
		Action: "register_proxy",
		Headers: map[string]string{
			"client-id": worker.client.Id,
			"tunnel-id": worker.tunnel.GetId(),
			"pub-conn-id": worker.publicConnId,
		},
	}

	protol.WriteMsg(conn, msg)

	// 启动代理管道
	worker.proxyConn = newProxyConn(conn)
	localConn, dialErr := net.Dial("tcp", worker.tunnel.ResolveAddress())

	if dialErr != nil {
		return dialErr
	}

	worker.localConn = localConn
	worker.proxyConn.pipe(worker.localConn)

	return nil
}

// Create one worker
func newWorker(client *Client, pubConnId string, tunnel tunnel.Tunnel) Worker{
	return &TcpWorker{
		client: client,
		publicConnId: pubConnId,
		tunnel: tunnel,
	}
}