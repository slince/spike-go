package conn

import (
	"container/list"
	"net"
	"sync"
)

type Pool struct {
	conns    *list.List
	max      int
	timeout  int
	lock     sync.Mutex
	cond     *sync.Cond
	callback func(pool *Pool)
}

func NewPool(max int, callback func(pool *Pool)) *Pool {
	return &Pool{
		conns:    list.New(),
		max:      max,
		lock:     sync.Mutex{},
		cond:     sync.NewCond(&sync.Mutex{}),
		callback: callback,
	}
}

func (p *Pool) Get() (conn net.Conn) {
	p.lock.Lock()
	var wait bool
	if p.conns.Len() == 0 {
		p.lock.Unlock()
		wait = true
		p.cond.L.Lock()
		p.callback(p)
		p.cond.Wait()
	}
	defer p.lock.Unlock()
	if wait {
		p.lock.Lock()
	}
	var ele = p.conns.Front()
	p.conns.Remove(ele)
	conn = ele.Value.(net.Conn)
	if wait {
		p.cond.L.Unlock()
	}
	return
}

func (p *Pool) Put(conn net.Conn) {
	p.lock.Lock()
	var wait = p.conns.Len() == 0
	p.conns.PushBack(conn)
	p.lock.Unlock()
	if wait {
		p.cond.Signal()
	}
}
