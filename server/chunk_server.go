package server

import (
	"errors"
	"fmt"
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/tunnel"
	"net"
)

type ChunkServer interface {
	// run this server
	run() error
	// get tunnel of chunk server
	getTunnel() tunnel.Tunnel
	// set proxy conn
	setProxyConn(pubConnId string, conn net.Conn) error
	// close the chunk server
	close()
}

// 监听公网接口
type TcpChunkServer struct {
	//对应的隧道
	tunnel *tunnel.TcpTunnel
	//服务的客户端
	client *Client
	//服务调度程序
	server *Server
	//公共请求
	pubConnsChan chan *PublicConn
	pubConns map[string]*PublicConn
	closeChan chan int
}

// 启动服务
func (chunkServer *TcpChunkServer) run() error {
	// enable listen
	listener, err := net.Listen("tcp", "127.0.0.1:" +
		chunkServer.tunnel.ServerPort)

	if err != nil {
		return errors.New("failed to create chunk server")
	}
	// listener accept
	go chunkServer.acceptConn(listener)
	// process public conn
	go chunkServer.processPublicConns()
	return nil
}

func (chunkServer *TcpChunkServer) acceptConn(listener net.Listener){
	for {
		select {
		case <- chunkServer.closeChan:
			listener.Close()
			break
		default:
			conn, err := listener.Accept()
			if err != nil {
				// handle error
				continue
			}
			//bytes := make([]byte, 1000)
			//conn.Read(bytes)
			//bytes,_ := ioutil.ReadAll(conn)
			//fmt.Println(string(bytes), "end")

			publicConn := newPublicConn(conn)
			chunkServer.pubConns[publicConn.Id] = publicConn
			chunkServer.pubConnsChan <- publicConn
		}
	}
}

func (chunkServer *TcpChunkServer) processPublicConns(){
	for {
		select {
		case <- chunkServer.closeChan:
			for _,pubConn := range chunkServer.pubConns {
				pubConn.close()
			}
			break
		case publicConn := <- chunkServer.pubConnsChan :
			//处理请求
			go chunkServer.handleConnection(publicConn)
		}
	}
}

// 处理公网请求
func (chunkServer *TcpChunkServer) handleConnection(pubConn *PublicConn) {
	//1.收到公网请求，请求客户端代理
	msg := protol.Protocol{
		Action: "request_proxy",
		Headers: map[string]string{
			"tunnel-id": chunkServer.tunnel.Id,
			"pub-conn-id": pubConn.Id,
		},
	}
	protol.WriteMsg(chunkServer.client.ctrlConn, &msg)

	// 2. 挂起当前公网请求
	//var proxyConn net.ctrlConn
	proxyConn := <- pubConn.proxyConnChan //从通道读取代理请求
	defer close(pubConn.proxyConnChan)

	fmt.Println("read a proxy conn")

	// 3. 管道请求
	pubConn.pipe(proxyConn)
	//delete(chunkServer.pubConnCollection, pubConn.Id)
}

// 获取隧道
func (chunkServer *TcpChunkServer) getTunnel() tunnel.Tunnel{
	return chunkServer.tunnel
}

// 设置代理链接
func (chunkServer *TcpChunkServer) setProxyConn(pubConnId string, conn net.Conn) error{
	fmt.Println("find a proxy conn")
	if pubConn, ok := chunkServer.pubConns[pubConnId];ok {
		pubConn.proxyConnChan <- conn
		return nil
	}
	return fmt.Errorf(`the public conn id "%s" is missing`, pubConnId)
}

// 关闭chunk server
func (chunkServer *TcpChunkServer) close() {
	chunkServer.closeChan <- 1
}

// http chunk server
type HttpChunkServer struct {
	TcpChunkServer
}