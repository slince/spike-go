package client

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type ProxyConn struct {
	conn net.Conn
	targetConn net.Conn
}

// close the proxy conn
func (proxyConn *ProxyConn) close() {
	proxyConn.conn.Close()
	if proxyConn.targetConn != nil {
		proxyConn.targetConn.Close()
	}
}

// 将当前请求管道输出到指定连接
func (proxyConn *ProxyConn) pipe(conn net.Conn) {

	defer conn.Close()
	defer proxyConn.conn.Close()

	var wait = new(sync.WaitGroup)
	wait.Add(2)

	go func() {
		for {
			io.Copy(conn, proxyConn.conn)
			fmt.Println("readed")
		}
		wait.Done()
	}()

	go func() {
		for {
			io.Copy(proxyConn.conn, conn)
			fmt.Println("copied")
		}
		wait.Done()
	}()
	wait.Wait()
}

// create new proxy conn
func newProxyConn(conn net.Conn) *ProxyConn{
	return &ProxyConn{
		conn: conn,
	}
}
