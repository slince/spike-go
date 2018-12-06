package client

import (
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type Worker interface {
	// Start the worker
	Start()
}

type TcpWorker struct {
	Client *Client
	publicConnId string
	localConn net.Conn
	proxyConn *ProxyConn
	tunnel tunnel.Tunnel
}

func (worker *TcpWorker) Start() error{
	conn, err := worker.createConnector()
	if err != nil {
		return err
	}
	// 发送消息给控制服务
	message := &protol.Protocol{
		Action: "register_proxy",
		Headers: map[string]string{
			"client-id": worker.Client.Id,
			"tunnel-id": worker.tunnel.GetId(),
			"public-connection-id": worker.publicConnId,
		},
	}
	conn.Write(message.ToBytes())

	reader := protol.NewReader(conn)
	for {
		messages,_ := reader.Read()
		for _, message := range messages {
			if message.Action == "start_proxy" { // 此时需要等待服务端传送start_proxy
				// 启动代理管道
				worker.proxyConn = &ProxyConn{
					Conn: conn,
				}
				localConn, dialErr := net.Dial("tcp", worker.tunnel.ResolveAddress())
				if dialErr != nil {
					worker.localConn = localConn
					worker.proxyConn.Pipe(worker.localConn)
				}
				break
			}
		}
	}

}

// 创建一个连接器连接控制服务器
func (worker *TcpWorker) createConnector() (net.Conn, error) {
	conn, err := net.Dial("tcp", worker.Client.ServerAddress)
	if err != nil {
		return conn, err
	}
	return conn, nil
}


