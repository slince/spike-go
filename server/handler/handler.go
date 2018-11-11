package handler

import (
	"github.com/slince/spike-go/protol"
	"github.com/slince/spike-go/server"
	"net"
)

type MessageHandler interface {
	// Handle the message
	Handle(message *protol.Protocol) error
}

type Handler struct{
	connection net.Conn
	server *server.Server
}