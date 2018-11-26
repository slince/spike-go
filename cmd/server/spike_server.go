package main

import (
	"github.com/slince/spike-go/server"
)

func main() {
	var ser *server.Server
	cfg,err := server.CreateConfigurationFromFile("./spiked.json")
	if err != nil {
		panic(err)
	}
	ser = server.NewServer(cfg)
	ser.Run()
}