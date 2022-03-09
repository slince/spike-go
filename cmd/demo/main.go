package main

import (
	"fmt"
	"net"
	"sync"
)

type Anima interface {
	Say()
}

type Dog struct {

}

func (d *Dog) Say(){

}

func isAnima(a Anima){
	fmt.Println(a)
}

var lock  sync.Mutex

func hello(){
	defer lock.Unlock()
	lock.Lock()
   	word()
}

func word(){
	defer lock.Unlock()
	lock.Lock()
	fmt.Println("hello world")
}

func main(){

	//http.ListenAndServe()
	//hello()


	var listen, err = net.Dial("tcp", "127.0.0.1:3306")

	if err != nil {
		panic(err)
	}

	for  {
		var buf = make([]byte, 16 * 1024)

		read, err := listen.Read(buf)

		if err != nil {
			panic(err)
		}

		fmt.Println(read, string(buf), "end")
	}
	//var reader io.Reader = bytes.NewReader([]byte("hello world i love china"))
	//var wantBytes = make([]byte, 4)
	//var reader io.Reader = bytes.NewReader([]byte("hello world i love china"))
	//lens, err := reader.Read(wantBytes)
	//fmt.Println(lens, err, string(wantBytes))
	//
	//wantBytes := make([]byte, 100)
	////lens, err := io.ReadFull(reader, wantBytes)
	//lens, err := reader.Read(wantBytes)
	//fmt.Println(lens, err, string(wantBytes))
	//
	//fmt.Println(protol.IntToBytes(123))
	//
	//var stopChan = make(chan int, 1)
	//
	//go (func() {
	//	stopChan <- 1
	//})()

	//for {
	//	select{
	//	case <- stopChan:
	//		println("stop")
	//		break
	//	default:
	//		//println("default")
	//	}
	//}
	println("end")
}