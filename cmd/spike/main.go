package main

import (
	"github.com/slince/spike/client"
)

func main(){
	var config, err = client.ConfigFromJsonFile("./spike.json")
	if err != nil {
		panic(err)
	}
	cli, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}
	err = cli.Start()
	if err != nil {
		panic(err)
	}
}
