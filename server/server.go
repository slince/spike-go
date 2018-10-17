package server

import "net"

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
func (server *Server) handleConnection(connection net.Conn) {

}


// Creates a new server.
func CreateServer(address string) *Server {
	return &Server{
		address,
		nil,
	}
}
