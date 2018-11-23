package main

import (
	"github.com/slince/spike-go/server"
)

func main() {
	var ser *server.Server
	ser = server.NewServer("0.0.0.0:8088", "./spike.log")
	ser.Run()
}