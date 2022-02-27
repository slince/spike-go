package conn

import (
	"fmt"
	"io"
	"net"
	"sync"
)

func copy(dst net.Conn, src net.Conn) (copied int64, err error, readErr bool) {
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
	if readErr {
		_ = src.Close()
	} else {
		_ = dst.Close()
	}
	return
}

func Combine(conn1 net.Conn, conn2 net.Conn, errCall func(alive net.Conn, err error)) (fromCopied int64, toCopied int64) {
	var wait sync.WaitGroup
	var pipe = func(conn1 net.Conn, conn2 net.Conn, copied *int64){
		defer wait.Done()
		var readErr bool
		var err error
		*copied, err, readErr = copy(conn2, conn1)
		fmt.Println("err combine:", readErr, err)
		if readErr {
			errCall(conn2, err)
		} else {
			errCall(conn1, err)
		}
	}
	wait.Add(2)
	go pipe(conn1, conn2, &fromCopied)
	go pipe(conn2, conn1, &toCopied)
	wait.Wait()
	return
}
