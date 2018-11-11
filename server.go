package main

import "github.com/slince/spike-go/server"

func main() {
	service := server.CreateServer("127.0.0.1:9090")
	service.Run()
}