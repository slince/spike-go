package client

import (
	"github.com/slince/spike/pkg/conn"
	"net"
	"sync"
)

type Pipe struct {
	src net.Conn
	dst net.Conn
	cond *sync.Cond
	factory func() (net.Conn,error)
}

func NewPipe(src net.Conn, factory func() (net.Conn,error)) *Pipe {
	return &Pipe{
		src: src,
		factory: factory,
	}
}

func (p *Pipe) getDst() net.Conn {
	if p.dst == nil {
		p.dst = p.factory()
	}
	return p.dst
}

func (p *Pipe) Combine(){
	var dstErr bool
	for {
		var dst = p.getDst()
		conn.Combine(p.src, dst, func(alive net.Conn, err error) {
			if alive == p.src {
				_ = dst.Close()
				p.dst = nil
				dstErr = true
			}
		})
		if dstErr {
			p.cond.L.Lock()
			p.cond.Wait()
			p.cond.L.Unlock()
		}
	}
}

func (p *Pipe) Signal(){
	p.cond.Signal()
}