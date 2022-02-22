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
}

type Factory struct {
	commands map[MsgType]reflect.Type
	types    map[reflect.Type]MsgType
}

/**
 * 注册支持的命令类型
 */
func (f *Factory) RegisterTypes(types map[MsgType]Command) {
	for msgType, command := range types {
		f.commands[msgType] = reflect.TypeOf(command)
		f.types[reflect.TypeOf(command)] = msgType
	}
}

/**
 * 将自定义命令序列化成message对象
 */
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

/**
 * 将message反序列化成command
 */
func (f *Factory) denormalize(msg Message) (command Command, err error) {
	t, ok := f.commands[msg.msgType]
	if !ok {
		err = errMsgType
		return
	}
	command = reflect.New(t).Interface().(Command)
	err = json.Unmarshal(msg.body, command)
	return
}

// NewFactory 工厂方法，创建新的命令工厂
func NewFactory() *Factory {
	return new(Factory)
}
