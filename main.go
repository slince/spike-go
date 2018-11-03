package main

import (
	"fmt"
	"time"
)


type Test struct {
	a int
	ch chan int
}

func (test Test) Run(){

	fmt.Println(<- test.ch)

	fmt.Println(test.a)
}


func main() {

	ch := make(chan int)

	test := Test{
		10,
		ch,
	}

	go test.Run()

	test.a = 20
	ch <- 10

	time.Sleep(10)
}