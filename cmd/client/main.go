package main

import (
	"github.com/slince/spike/client"
)

func main(){
	cli := client.NewClient("127.0.0.1", 8808, "admin", "admin")
	err := cli.Start()
	if err != nil {
		panic(err)
	}
}
