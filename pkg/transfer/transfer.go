package transfer

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Bridge struct {
	reader io.Reader
	writer io.Writer
	ft     *Factory
}

var ft = NewFactory()
var parser = NewParser()

/**
 * 写入命令到该通道
 */
func (b *Bridge) Write(command Command) error {
	msg, err := ft.normalize(command)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.BigEndian, msg.length)
	if err != nil {
		return err
	}

	buf.Write(msg.body)
	return nil
}

/**
 * 从通道读取一个command
 */
func (b *Bridge) Read() (command Command, err error) {
	msg, err := parser.parse(b.reader)
	if err != nil {
		return
	}

	command, err = ft.denormalize(msg)
	return
}

/**
 * 创建通道
 */
func NewBridge(reader io.Reader, writer io.Writer) *Bridge {
	return &Bridge{reader, writer, ft}
}
