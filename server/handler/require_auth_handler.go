package handler

import (
	"github.com/slince/jinbox/protol"
	"github.com/slince/jinbox/server"
)

type RequireAuthHandler struct {
	Handler
	client *server.Client
}

func (hd *RequireAuthHandler) isAuthenticated(message *protol.Protocol) bool{
	clientId, ok := message.Headers["client-id"]
	if !ok {
		return false
	}
	if client, ok := hd.server.Clients[clientId]; ok{
		hd.client = client
		return true
	}
	return false
}