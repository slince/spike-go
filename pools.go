package main

import (
	"container/list"
	"net"
	"sync"
	"time"
)

//Conn - Wrap of net.Conn
type Conn struct {
	conn     net.Conn
	deadline time.Time
	timeout  int
}

//Read - Compatible io.Reader
func (s *Conn) Read(p []byte) (int, error) {
	n, err := s.conn.Read(p)
	if err != nil {
		return n, err
	}
	s.setDeadline()
	return n, nil
}

//Write - Compatible io.Write
func (s *Conn) Write(p []byte) (int, error) {
	n, err := s.conn.Write(p)
	if err != nil {
		return n, err
	}
	s.setDeadline()
	return n, nil
}

//Close - Close the connection
func (s *Conn) Close() error {
	err := s.conn.Close()
	s.deadline = time.Now()
	return err
}

func (s *Conn) setDeadline() {
	s.deadline = time.Now().Add(time.Duration(s.timeout) * time.Second)
	s.conn.SetDeadline(s.deadline)
}

//Timeout - Test the connection is timeout. return true/false.
func (s *Conn) Timeout() bool {
	return time.Now().Unix() >= s.deadline.Unix()
}

//Pool - Please create Pool with function NewPool
type Pool struct {
	lock        *sync.Mutex
	cond        *sync.Cond
	max         int
	timeout     int
	newFunc     func() (net.Conn, error)
	connections *list.List
}

func NewPool(max, timeout int, factory func() (net.Conn, error)) *Pool {
	l := new(sync.Mutex)
	return &Pool{lock: new(sync.Mutex), cond: sync.NewCond(l), max: max, timeout: timeout, newFunc: factory,
		connections: list.New()}
}

func (s *Pool) newConn() (*Conn, error) {
	conn, err := s.newFunc()
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(time.Duration(s.timeout) * time.Second)
	conn.SetDeadline(deadline)
	return &Conn{conn: conn, deadline: deadline, timeout: s.timeout}, nil
}

//Get - Get a new connection from the Pool
func (s *Pool) Get() (*Conn, error) {
	var res *Conn
	var wait bool
	s.lock.Lock()
	if s.max > 0 {
		s.max--
	} else {
		wait = true
	}
	s.lock.Unlock()

	if wait {
		s.cond.L.Lock()
		s.cond.Wait()
		s.cond.L.Unlock()
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.connections.Len() > 0 {
		e := s.connections.Front()
		res = e.Value.(*Conn)
		s.connections.Remove(e)
		if res.Timeout() {
			go res.Close()
		} else {
			return res, nil
		}
	}

	return s.newConn()

}

//Put - Put a connection back to the Pool
func (s *Pool) Put(conn1 *Conn) {
	var wait bool
	s.lock.Lock()
	s.connections.PushBack(conn1)
	if s.max == 0 {
		wait = true
	}
	s.lock.Unlock()
	if wait {
		s.cond.Signal()
	}
}

//Close - Close all connections in this Pool
func (s *Pool) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	e := s.connections.Front()
	if e.Value != nil {
		conn := e.Value.(*Conn)
		conn.Close()
		e = e.Next()
	}
	s.connections.Init()
}
