package conn

import (
	"io"
	"net"
	"sync"
)

func copy(dst io.Writer, src io.Reader) (copied int64, err error, readErr bool) {
	var buf = make([]byte, 32 * 1024)
	for {
		read, err1 := src.Read(buf)
		if read > 0 {
			write, err2 := dst.Write(buf[0:read])
			copied += int64(write)
			if err2 != nil {
				err = err2
				readErr = false
				break
			}
			if read != write {
				err = io.ErrShortWrite
				readErr = false
				break
			}
		}
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			readErr = true
			break
		}
	}
	return
}

func Combine(conn1 net.Conn, conn2 net.Conn, errCallback func(con net.Conn)) (fromCopied int64, toCopied int64) {
	var wait sync.WaitGroup
	var pipe = func(conn1 net.Conn, conn2 net.Conn, copied *int64){
		defer wait.Done()
		var readErr bool
		*copied, _, readErr = copy(conn2, conn1)
		if readErr {
			errCallback(conn1)
		} else {
			errCallback(conn2)
		}
	}
	go pipe(conn1, conn2, &fromCopied)
	go pipe(conn2, conn1, &toCopied)
	wait.Add(2)
	wait.Wait()
	return
}
