package main

import (
	"fmt"
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

	//reader := strings.NewReader("哈哈哈哈")
	//
	//reader2 := bufio.NewReader(reader)
	//
	//bytes := make([]byte, 50)
	//
	//readedBytes, _ := reader2.Read(bytes)
	//
	//fmt.Println(readedBytes)
	//fmt.Println(string(bytes))

	//
	//var map1  = map[string]string{
	//	"a": "asda",
	//	"b": "asda",
	//}
	//
	//func(map2 map[string]string){
	//
	//	(map2)["c"] = "12312"
	//
	//}(map1)
	//
	//fmt.Println(map1)


	var sl []int
	fmt.Println(sl == nil)
	sl = append(sl, 10)
	sl = append(sl, 12)
	fmt.Println(sl)
}