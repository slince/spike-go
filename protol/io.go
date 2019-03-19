package protol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)


import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Protocol Reader
type Reader struct {
	reader io.Reader
}

// read a message from io.reader
func (reader *Reader) Read() (*Protocol, error){
	var length int64
	err := binary.Read(reader.reader, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	fmt.Println(length)

	readBytes, err := ioutil.ReadAll(reader.reader)
	if err != nil{
		return nil, err
	}
	protocol := &Protocol{}
	err = json.Unmarshal(readBytes, protocol)

	if err != nil {
		return nil, err
	}

	return protocol, nil
}


// Create a new protocol reader
func NewReader(reader io.Reader) *Reader{
	return &Reader{
		reader,
	}
}


type Writer struct {
	writer io.Writer
}

func (writer *Writer)Write(protocol *Protocol) (int,error){
	con := protocol.ToBytes()
	leng := len(con)
	var msg = append(IntToBytes(leng), con...)
	return writer.writer.Write(msg)
}

// Create a new protocol writer
func NewWriter(writer io.Writer) *Writer{
	return &Writer{
		writer,
	}
}

func IntToBytes(num int) []byte{
	x := int64(num)
	var buf = bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

type IO struct {
	conn net.Conn
	reader *Reader
	writer *Writer
}
func (io *IO) Write(protocol *Protocol) (int, error){
	return io.writer.Write(protocol)
}

func (io *IO) Read() (*Protocol, error){
	return io.reader.Read()
}

func NewIO(conn net.Conn) *IO{
	return &IO{
		conn,
		NewReader(conn),
		NewWriter(conn),
	}
}




