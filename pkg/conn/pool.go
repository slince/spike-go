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
	callback func()
}

func NewPool(max int, callback func()) *Pool {
	return &Pool{
		conns:    list.New(),
		max:      max,
		lock:     sync.Mutex{},
		cond:     sync.NewCond(&sync.Mutex{}),
		callback: callback,
	}
}

func (p *Pool) Get() (conn net.Conn) {
	defer p.lock.Unlock()
	p.lock.Lock()
	var wait = false
	if p.conns.Len() == 0 {
		wait = true
		p.lock.Unlock()
		p.cond.L.Lock()
		p.callback()
		p.cond.Wait()
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
	var wait = p.conns.Len() == 0
	p.lock.Lock()
	p.conns.PushBack(conn)
	p.lock.Unlock()
	if wait {
		p.cond.Signal()
	}
}
