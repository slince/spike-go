package main

import (
	"bufio"
	"fmt"
	"strings"
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

	reader := strings.NewReader("哈哈哈哈")

	reader2 := bufio.NewReader(reader)

	bytes := make([]byte, 50)

	readedBytes, _ := reader2.Read(bytes)

	fmt.Println(readedBytes)
	fmt.Println(string(bytes))
}