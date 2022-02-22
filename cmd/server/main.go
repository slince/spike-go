package main

import (
	"github.com/slince/spike/pkg/auth"
	"github.com/slince/spike/server"
)

func main(){
	var users = make([]*auth.GenericUser, 2)
	users[0] = &auth.GenericUser{
		Username: "admin",
		Password: "admin",
	}
	users[1] = &auth.GenericUser{
		Username: "test",
		Password: "123456",
	}
	var au = auth.NewSimpleAuth(users)
	ser := server.NewServer("127.0.0.1", 8808, au)
	err := ser.Start()
	if err != nil {
		panic(err)
	}
}
