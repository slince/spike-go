package handler

import (
	"fmt"
	"github.com/slince/spike-go/protol"
)

type RegisterProxyHandler struct{
	RequireAuthHandler
}

func (hd *RegisterProxyHandler) Handle(message *protol.Protocol){
	tunnelId, ok := message.Headers["tunnel-id"]
	if !ok {
		hd.connection.Write([]byte("missing tunnel id"))
		hd.connection.Close()
		return
	}
	chunkServer := hd.server.FindChunkServer(tunnelId)
	if chunkServer == nil {
		hd.connection.Write([]byte(fmt.Sprintf("the chunk server %s is not found", tunnelId)))
		hd.connection.Close()
		return
	}
	publicConnectionId,ok := message.Headers["public-connection-id"]
	if !ok { //错误的注册代理协议
		hd.connection.Write([]byte("missing public id"))
		hd.connection.Close()
	}
	// set proxy connection
	chunkServer.SetProxyConnection(publicConnectionId, hd.connection)
}