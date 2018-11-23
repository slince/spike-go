package main

import (
	"github.com/slince/spike-go/server"
)

func main() {
	var ser *server.Server
	ser = server.NewServer("127.0.0.1:8090", "./spike.log")
	ser.Run()
}