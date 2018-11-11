package main

import (
	"github.com/slince/spike-go/server"
)

func main() {

	var s server.Server

	s = server.NewServer("0.0.0.0:8088")

	s.Run()
}