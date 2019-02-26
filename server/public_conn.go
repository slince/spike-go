package server

import (
	"fmt"
	"github.com/rs/xid"
	"io"
	"net"
	"sync"
)

// 公网请求
type PublicConn struct {
	Id string
	Conn net.Conn
	ProxyConnChan chan net.Conn
	pubLock sync.RWMutex
	proxyLock sync.RWMutex
}

func (pubConn *PublicConn) Pipe(proxyConn net.Conn) {

	defer proxyConn.Close()
	defer pubConn.Conn.Close()

	var wait sync.WaitGroup
	wait.Add(2)
	go func() { // 从公网请求读数据并写入到代理请求
		for {
			fmt.Println("pub read start")
			pubConn.pubLock.RLock() // 公众请求读锁
			pubConn.proxyLock.Lock() // 代理请求写锁
			fmt.Println("pub copy start")
			writtenBytes, err := io.Copy(proxyConn, pubConn.Conn)
			fmt.Println("pub copied")
			if err != nil { //读取出错两者都关闭
				panic(err)
				proxyConn.Close()
				pubConn.Conn.Close()
				return
			}
			fmt.Println("pub read bytes:", writtenBytes)
			pubConn.pubLock.RUnlock()
			pubConn.proxyLock.Unlock()
		}
	}()
	go func() { // 从代理请求读取并写入到公众请求
		for {
			fmt.Println("proxy read start")
			pubConn.proxyLock.RLock()
			pubConn.pubLock.Lock()
			fmt.Println("proxy copy start")
			writtenBytes, err := io.Copy(pubConn.Conn, proxyConn)
			fmt.Println("proxy copied")
			if err != nil { //读取出错两者都关闭
				panic(err)
				proxyConn.Close()
				pubConn.Conn.Close()
				return
			}
			fmt.Println("proxy read bytes:", writtenBytes)
			pubConn.proxyLock.RUnlock()
			pubConn.pubLock.Unlock()
		}
	}()

	wait.Wait()
}

// Create a public connection.
func NewPublicConn(conn net.Conn) *PublicConn {
	ch := make(chan net.Conn)
	return &PublicConn{
		Id: xid.New().String(),
		Conn: conn,
		ProxyConnChan: ch,
	}
}