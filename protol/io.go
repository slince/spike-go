package protol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	errorBadLength = errors.New("bad protocol length")
)

// Protocol Reader
type Reader struct {
	reader io.Reader
}

// read a message from io.reader
func (reader *Reader) Read() (protocol *Protocol, err error){
	// read length bytes
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(reader.reader, lenBuf)

	if err != nil {
		fmt.Println(err)
		err = errorBadLength
		return
	}
	var length = BytesToInt(lenBuf)

	var conBuf = make([]byte, length)
	_, err = io.ReadFull(reader.reader, conBuf)

	if err != nil{
		return
	}
	protocol = &Protocol{}
	err = json.Unmarshal(conBuf, protocol)
	return
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
	x := uint32(num)
	var buf = bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

func BytesToInt(buf []byte) int{
	var x uint32
	binary.Read(bytes.NewReader(buf), binary.BigEndian, &x)
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

func ReadMsg(conn net.Conn){
	NewIO(conn).Read()
}

func WriteMsg(conn net.Conn, msg *Protocol){
	NewIO(conn).Write(msg)
}




