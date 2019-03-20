package client

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type ProxyConn struct {
	conn net.Conn
}

// close the proxy conn
func (proxyConn *ProxyConn) close() {
	proxyConn.conn.Close()
}

// 将当前请求管道输出到指定连接
func (proxyConn *ProxyConn) pipe(localConn net.Conn) {

	var wait = new(sync.WaitGroup)
	wait.Add(2)

	go func() {
		for {
			_,err := io.Copy(localConn, proxyConn.conn)
			if err != nil {
				localConn.Close()
				fmt.Println("proxy closed")
				break
			}
		}
		wait.Done()
	}()

	go func() {
		for {
			_,err := io.Copy(proxyConn.conn, localConn)
			if err != nil {
				proxyConn.close()
				fmt.Println("local closed")
				break
			}
		}
		wait.Done()
	}()
	wait.Wait()
	fmt.Println("proxy end")
}

// create new proxy conn
func newProxyConn(conn net.Conn) *ProxyConn{
	return &ProxyConn{
		conn: conn,
	}
}
