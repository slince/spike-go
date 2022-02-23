package main

import (
	"github.com/slince/spike/server"
)

func main(){
	var config, err = server.ConfigFromJsonFile("./spiked.json")
	if err != nil {
		panic(err)
	}
	ser, err := server.NewServer(config)
	if err != nil {
		panic(err)
	}
	err = ser.Start()
	if err != nil {
		panic(err)
	}
}
