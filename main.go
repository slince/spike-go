package main

func main(){

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

	var stopChan = make(chan int, 1)

	go (func() {
		stopChan <- 1
	})()

	for {
		select{
		case <- stopChan:
			println("stop")
			break
		default:
			//println("default")
		}
	}
	println("end")
}