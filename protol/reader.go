package protol

import (
	"bytes"
	"errors"
	"io"
)

// Protocol Reader
type Reader struct {
	rd io.Reader
	incoming *bytes.Buffer
	endCharacter string
}

// Read a message from io.reader
func (reader *Reader) ReadMessage() (*Protocol, error) {

	for {
		incoming, err := reader.read(50)

		if err == nil && len(incoming) > 0 {
			//
			for len(incoming) > 0 {

				if reader.endCharacter == "" {
					firstCharacter := string(incoming[0])

					if firstCharacter == "[" {
						reader.endCharacter = "]"
					} else if firstCharacter == "{" {
						reader.endCharacter = "}"
					} else {
						return nil, errors.New("bad json start character")
					}

				}
				pos := bytes.IndexAny(incoming, reader.endCharacter)

				if pos == -1 {
					break
				}



			}
		}
	}

}

// read
func (reader *Reader)read(n int) ([]byte, error) {
	incoming := make([]byte, n)

	n, err := reader.rd.Read(incoming)

	if err == nil && n > 0 {
		reader.incoming.Write(incoming)
	}

	return incoming, nil
}


// Create a new protocol reader
func NewReader(rd io.Reader) *Reader{
	return &Reader{
		rd,
		bytes.NewBuffer(make([]byte, 0, 50)),
		"",
	}
}