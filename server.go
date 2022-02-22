package main

import "github.com/slince/spike/server"

func main() {
	service := server.CreateServer("127.0.0.1:9090")
	service.Run()
}