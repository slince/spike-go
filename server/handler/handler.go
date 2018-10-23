package handler

import (
	"github.com/slince/jinbox/protol"
	"github.com/slince/jinbox/server"
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