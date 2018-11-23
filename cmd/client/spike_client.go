package main

import (
	"github.com/slince/spike-go/client"
)

func main() {

	var clt *client.Client
	var cfg *client.Configuration

	cfg,_ = client.CreateConfigurationFromFile("./spike.json")
	clt = client.NewClient(cfg)
	clt.Start()
}