package conn

import (
	"io"
	"net"
	"sync"
)

type Pipe struct {
	src        io.Reader
	dst        io.Writer
	readError  func(src io.Reader)
	writeError func(dst io.Writer)
}

func NewPipe(src io.Reader, dst io.Writer, readError func(src io.Reader), writeError func(dst io.Writer)) *Pipe {
	return &Pipe{
		src:        src,
		dst:        dst,
		readError:  readError,
		writeError: writeError,
	}
}

func (p *Pipe) Pipe() {
	for {
		var buf = make([]byte, 10)
		var _, err = p.src.Read(buf)
		if err != nil {
			_, err = p.dst.Write(buf)
			if err != nil {
				continue
			} else {
				p.writeError(p.dst)
			}
		} else {
			p.readError(p.src)
		}
		break
	}
}

func Combine(conn1 net.Conn, conn2 net.Conn, readError func(src io.Reader), writeError func(dst io.Writer)) {
	var wait sync.WaitGroup
	wait.Add(2)
	go (func() {
		defer wait.Done()
		//var pipe = NewPipe(conn1, conn2, readError, writeError)
		//pipe.Pipe()
		io.Copy(conn2, conn1)
	})()
	go (func() {
		defer wait.Done()
		var pipe = NewPipe(conn2, conn1, readError, writeError)
		pipe.Pipe()
	})()
	wait.Wait()
}
