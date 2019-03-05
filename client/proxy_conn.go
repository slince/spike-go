package client

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type ProxyConn struct {
	Conn net.Conn
}

// 将当前请求管道输出到指定连接
func (proxyConn *ProxyConn) Pipe(conn net.Conn) {

	defer conn.Close()
	defer proxyConn.Conn.Close()

	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		io.Copy(conn, proxyConn.Conn)
		fmt.Println("readed")
		wait.Done()
	}()
	go func() {
		io.Copy(proxyConn.Conn, conn)
		fmt.Println("copied")
		wait.Done()
	}()

	wait.Wait()
}
