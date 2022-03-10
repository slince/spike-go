package conn

import (
	"errors"
	"net"
	"time"
)

type Pool struct {
	num int
	conns    chan net.Conn
	timeout  int
	callback func(pool *Pool)
}

func NewPool(max int, timeout int, callback func(pool *Pool)) *Pool {
	return &Pool{
		conns:    make(chan net.Conn, max),
		timeout: timeout,
		callback: callback,
	}
}

func (p *Pool) Get() (net.Conn, error) {
	if len(p.conns) < cap(p.conns) {
		p.callback(p)
	}
	for {
		var after = time.After(time.Second * time.Duration(p.timeout))
		select {
			case <- after:
				return nil, errors.New("timeout to get conn")
			case con := <- p.conns:
				return con, nil
		}
	}
}

func (p *Pool) Put(conn net.Conn) {
	p.conns <- conn
}
