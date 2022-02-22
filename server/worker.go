package server

import (
	"github.com/slince/spike/pkg/tunnel"
	"net"
	"strconv"
)

type Worker struct {
	tun *tunnel.Tunnel
	conn *net.Conn
}


func (w *Worker) Start() (err error){
	address := "127.0.0.1:" + strconv.Itoa(int(w.tun.ServerPort))
	socket, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			return
		}
		go w.handleConn(conn)
	}
}


func (w *Worker)handleConn(conn net.Conn){

}
