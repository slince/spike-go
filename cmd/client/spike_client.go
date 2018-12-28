package main

import (
	"github.com/slince/spike-go/client"
)

func main() {

	var clt *client.Client
	cfg,err := client.CreateConfigs("./spike.json")

	if err != nil {
		panic(err)
	}

	clt = client.NewClient(cfg)
	clt.Start()
}