package transfer

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	errMsgType      = errors.New("message type error")
	errMaxMsgLength = errors.New("message length exceed the limit")
	//errMsgLength    = errors.New("message length error")
	errMsgFormat = errors.New("message format error")
)

type Command interface {
	setRawBody(body  []byte)
}

type BaseCommand struct {
	rawBody []byte
}

func (c *BaseCommand) String() string{
	return string(c.rawBody)
}
func (c *BaseCommand) setRawBody(body []byte){
	c.rawBody = body
}

type Factory struct {
	commands map[MsgType]reflect.Type
	types    map[reflect.Type]MsgType
}

func (f *Factory) RegisterTypes(types map[MsgType]Command) {
	for msgType, command := range types {
		f.commands[msgType] = reflect.TypeOf(command)
		f.types[reflect.TypeOf(command)] = msgType
	}
}

func (f *Factory) normalize(command Command) (msg Message, err error) {
	msgType, ok := f.types[reflect.TypeOf(command)]
	if !ok {
		err = errMsgType
		return
	}
	body, err := json.Marshal(command)
	if err != nil {
		return
	}
	msg = newMessage(msgType, body)
	return
}

func (f *Factory) denormalize(msg Message) (command Command, err error) {
	t, ok := f.commands[msg.msgType]
	if !ok {
		err = errMsgType
		return
	}
	command = reflect.New(t).Interface().(Command)
	err = json.Unmarshal(msg.body, command)
	command.setRawBody(msg.body)
	return
}

func NewFactory() *Factory {
	return &Factory{
		make(map[MsgType]reflect.Type, 0),
		make(map[reflect.Type]MsgType, 0),
	}
}
