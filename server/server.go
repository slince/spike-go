package server

import (
	"bufio"
	"github.com/slince/jinbox/event"
	"github.com/slince/jinbox/protol"
	"net"
)

type Server struct {
	Address string
	socket net.Listener
	dispatcher *event.Dispatcher
}

// Run the server
func (server *Server) Run() {

	// register all listeners
	server.registerListeners()

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

// Register all listeners
func (server *Server)registerListeners() {
	var callback = func(event *event.Event) {

	}
	server.dispatcher.On(serverInit, listener)
}

// handle connection from client.
func (server *Server) handleConnection(connection net.Conn) error{
	str, err := bufio.NewReader(connection).ReadString('\n')

	if err != nil {
		return err
	}

	message,err := protol.FromJsonString(str)

	if err != nil {
		return err
	}

	ev := event.NewEvent("message", map[string]interface{}{"message":  message})

	server.dispatcher.Fire(ev)

	return nil
}


// Creates a new server.
func CreateServer(address string) *Server {
	return &Server{
		address,
		nil,
		event.NewDispatcher(),
	}
}
