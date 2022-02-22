package transfer

import (
	"io"
)

type Bridge struct {
	reader io.Reader
	writer io.Writer
	ft     *Factory
}

var parser = NewParser()

/**
 * 写入命令到该通道
 */
func (b *Bridge) Write(command Command) error {
	msg, err := b.ft.normalize(command)
	if err != nil {
		return err
	}
	buf, err := parser.pack(msg)
	if err != nil {
		return err
	}
	_, err = b.writer.Write(buf)
	return err
}

/**
 * 从通道读取一个command
 */
func (b *Bridge) Read() (command Command, err error) {
	msg, err := parser.parse(b.reader)
	if err != nil {
		return
	}
	command, err = b.ft.denormalize(msg)
	return
}

func (b *Bridge) Supports(types map[MsgType]Command) {
	b.ft.RegisterTypes(types)
}

func NewBridge(ft *Factory, reader io.Reader, writer io.Writer) *Bridge {
	return &Bridge{reader, writer, ft}
}
