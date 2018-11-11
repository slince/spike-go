package util

import (
	"io"
	"os/exec"
)


func Pipe(from io.Reader, to io.Writer) error{

	cmd1 := exec.Command("ps", "aux")

	cmd1.StdoutPipe()

	var  (
		rChannel chan []byte
	)
	go func() {
		for {
			bytes := make([]byte, 50)
			_, err := from.Read(bytes)
			if err != nil {

			}
			rChannel <- bytes
		}
	}()

	go func() {
		for {
			data := <- rChannel
			to.Write(data)
		}
	}()
}
