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
	conn net.Conn
	proxyConnChan chan net.Conn
}

func (pubConn *PublicConn) pipe(proxyConn net.Conn) {

	defer proxyConn.Close()
	defer pubConn.conn.Close()

	var wait = sync.WaitGroup{}
	wait.Add(2)

	go func() { // 从代理请求读取并写入到公众请求
		for {
			_,err := io.Copy(pubConn.conn, proxyConn)
			if err != nil {
				pubConn.close()
				fmt.Println("proxy closed")
				break
			}
		}
		wait.Done()
	}()

	go func() {  //从公网请求读数据并写入到代理请求
		for {
			_,err := io.Copy(proxyConn, pubConn.conn)
			if err != nil { //读取出错两者都关闭
				proxyConn.Close()
				fmt.Println("public closed")
				break
			}
		}
		wait.Done()
	}()

	wait.Wait()
	fmt.Println("pub end")
}

// close
func (pubConn *PublicConn) close() {
	pubConn.conn.Close()
}

// Create a public conn.
func newPublicConn(conn net.Conn) *PublicConn {
	ch := make(chan net.Conn)
	return &PublicConn{
		Id: xid.New().String(),
		conn: conn,
		proxyConnChan: ch,
	}
}