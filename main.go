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

	//
	//var sl []int
	//fmt.Println(sl == nil)
	//sl = append(sl, 10)
	//sl = append(sl, 12)
	//fmt.Println(sl)

	//
	//type Bar struct {
	//	Bar1 string `json:"bar_1"`
	//	Bar2 string `json:"bar_2"`
	//}
	//
	//type Foo struct {
	//	Bar interface{}
	//}
	//
	//foo := Foo{
	//	Bar: Bar{
	//		Bar1: "tao",
	//		Bar2: "sikai",
	//	},
	//}
	//jsonS, _ := json.Marshal(foo)
	//fmt.Println(string(jsonS))
	//
	//foo1 := &Foo{
	//	Bar: Bar{
	//
	//	},
	//}
	//json.Unmarshal(jsonS, foo1)
	//fmt.Print(foo1.Bar)

	//
	//var map1 = map[string]string{
	//	"foo": "bar",
	//	"bar": "baz",
	//}
	//
	//bytes.NewBufferString("").String()
	//fmt.Println(map1["fooz"])


	var sl []int
	sl = append(sl, 10)
	//sl[0] = 10
	fmt.Println(sl, len(sl))
}