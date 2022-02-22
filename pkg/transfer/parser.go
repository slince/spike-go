package transfer

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	MaxMsgLength = 10240
)

type MsgType uint8

type Message struct {
	msgType MsgType
	length  uint64
	body    []byte
}

func newMessage(msgType MsgType, body []byte) Message {
	return Message{
		msgType: msgType,
		body:    body,
		length:  uint64(len(body)),
	}
}

type Parser struct {
}

func (p *Parser) parse(r io.Reader) (msg Message, err error) {
	msg = Message{}
	err = p.meta(r, msg)
	if err != nil {
		return
	}

	body := make([]byte, msg.length)
	n, err := io.ReadFull(r, body)
	if err != nil {
		return
	}

	if uint64(n) != msg.length {
		err = errMsgFormat
		return
	}
	msg.body = body
	return
}

func (p *Parser) pack(msg Message) (buffer []byte, err error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(byte(msg.msgType))

	err = binary.Write(buf, binary.BigEndian, msg.length)
	if err != nil {
		return
	}

	buf.Write(msg.body)
	buffer = buf.Bytes()
	return
}

func (p *Parser) meta(r io.Reader, msg Message) error {
	buffer := make([]byte, 1)
	_, err := r.Read(buffer)

	if err != nil {
		return err
	}

	msg.msgType = MsgType(buffer[0])
	//var bodyLength uint64
	err = binary.Read(r, binary.BigEndian, &msg.length)
	if err != nil {
		return err
	}
	if msg.length > MaxMsgLength {
		return errMaxMsgLength
	}
	return nil
}

func newParser() *Parser {
	return new(Parser)
}
