package server

import (
	"bufio"
	"github.com/slince/jinbox/protol"
	"net"
)

type Server struct {
	Address string
	socket net.Listener
}

// Run the server
func (server *Server) Run() {
	var err error
	server.socket ,err = net.Listen("tcp", server.Address)
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := server.socket.Accept()
		if err != nil {
			// handle error
			continue
		}
		go server.handleConnection(conn)
	}
}

// handle connection from client.
func (server *Server) handleConnection(connection net.Conn) error{
	str, err := bufio.NewReader(connection).ReadString('\n')

	if err != nil {
		return err
	}

	protocol,err := protol.FromJsonString(str)

	if err != nil {
		return err
	}

	
}


// Creates a new server.
func CreateServer(address string) *Server {
	return &Server{
		address,
		nil,
	}
}
