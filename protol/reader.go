package protol

import (
	"bytes"
	"errors"
	"io"
)

const maxReadLength = 60

// Protocol Reader
type Reader struct {
	rd io.Reader
	Incoming *bytes.Buffer
	endCharacter string
}

// Read a message from io.reader
func (reader *Reader) Read() ([]*Protocol, error) {
	var messages = make([]*Protocol, 0, 5)

	for {
		chunk := make([]byte, maxReadLength)
		read, err := reader.rd.Read(chunk)

		if err == nil {

			if read < maxReadLength {
				chunk = chunk[:read]
			}

			for len(chunk) > 0 {

				if reader.endCharacter == "" {
					firstCharacter := string(chunk[0])

					if firstCharacter == "[" {
						reader.endCharacter = "]"
					} else if firstCharacter == "{" {
						reader.endCharacter = "}"
					} else {
						return nil, errors.New("bad json start character:" + firstCharacter)
					}
				}

				pos := bytes.Index(chunk, []byte(reader.endCharacter))

				if pos == -1 {
					reader.Incoming.Write(chunk)
					break
				}

				reader.Incoming.Write(chunk[:pos+1])
				chunk = chunk[pos+1:]

				message, err := FromJsonString(reader.Incoming.String())

				if err == nil {
					messages = append(messages, message)
					reader.Incoming.Reset()
					reader.endCharacter = ""
				}
			}

			if read < maxReadLength {
				break
			}

		} else {
			break
		}
	}
	return messages, nil
}


// Create a new protocol reader
func NewReader(rd io.Reader) *Reader{
	return &Reader{
		rd,
		bytes.NewBuffer(make([]byte, 0, 50)),
		"",
	}
}