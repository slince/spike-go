package protol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
)

var (
	errorReadError = errors.New("read error")
)

// Protocol Reader
type Reader struct {
	reader io.Reader
}

// read a message from io.reader
func (reader *Reader) Read() (protocol *Protocol, err error){
	var length int64
	err = binary.Read(reader.reader, binary.BigEndian, &length)
	if err != nil {
		return
	}

	var readBytes = make([]byte, length)
	readLength, err := io.ReadFull(reader.reader, readBytes)

	if err != nil{
		return
	}
	protocol = &Protocol{}
	err = json.Unmarshal(readBytes, protocol)

	if err != nil {
		return
	}

	if int64(readLength) != length {
		err = errorReadError
		return
	}
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
	x := int64(num)
	var buf = bytes.NewBuffer(make([]byte, 0))
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
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




