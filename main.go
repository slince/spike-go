package main

import "fmt"

func main() {

	//var str = "woshizhongguoren"
	//
	//var bytesBuf = []byte(str)
	//
	//fmt.Println(strings.Index(str, "z"))
	//fmt.Println(bytes.Index(bytesBuf, []byte("z")))
	//
	//strArr := []int{1,2,3,4}
	//
	//fmt.Println(strArr[0:3])
	//fmt.Println(strArr[1:3])
	//
	//arr := make([]byte, 20)
	//
	//buf := bytes.NewBufferString("imtaosikai")
	//
	//n,err := buf.Read(arr)
	//
	//fmt.Println(n)
	//fmt.Println(err)
	//fmt.Println(arr)
	//fmt.Println(len(arr))
	//fmt.Println(cap(arr))


	arr := []int{1,2,3,4,5}

	fmt.Println(arr[:3])
	fmt.Println(arr[:2])
	fmt.Println(arr[1:3])
	fmt.Println(arr[1:2])
}