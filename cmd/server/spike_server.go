package main

import "github.com/slince/spike-go/server"

func main() {
	err := server.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}