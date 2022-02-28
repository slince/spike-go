package conn

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func copy(dst net.Conn, src net.Conn, stop chan bool) (copied int64, err error, readErr bool) {
	var buf = make([]byte, 32 * 1024)
	readErr = true
	Handle:
	for {
		select {
		case <- stop:
			fmt.Println("chan stop")
			break Handle
		default:
			err = src.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				readErr = true
				break Handle
			}
			read, err1 := src.Read(buf)
			if os.IsTimeout(err1) {
				fmt.Println("超时检查。。。", GetGID())
				break
			}
			if read > 0 {
				write, err2 := dst.Write(buf[0:read])
				copied += int64(write)
				if err2 != nil {
					err = err2
					readErr = false
					break Handle
				}
				if read != write {
					err = io.ErrShortWrite
					readErr = false
					break Handle
				}
			}
			if err1 != nil {
				if err1 != io.EOF {
					err = err1
				}
				err = err1
				readErr = true
				break Handle
			}
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
	var stop = make(chan bool, 1)
	var pipe = func(conn1 net.Conn, conn2 net.Conn, copied *int64){
		defer wait.Done()
		var readErr bool
		var err error
		*copied, err, readErr = copy(conn2, conn1, stop)
		fmt.Println("err combine:", readErr, err)
		stop <- true
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
	close(stop)
	return
}
