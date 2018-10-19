package protol

import (
	"bytes"
	"testing"
)

const jsonMessage = `{"action":"register","body":{"error":""},"headers":{"status":"200"}}`

func TestReadFullJson(t *testing.T) {

	reader := NewReader(bytes.NewBuffer([]byte(jsonMessage)))

	messages,err := reader.Read()

	if err != nil {
		t.Errorf(err.Error())
	}

	if messages[0].Action != "register" {
		t.Errorf("error parse action")
	}
}

func TestReadManyJson(t *testing.T) {

	reader := NewReader(bytes.NewBuffer([]byte(jsonMessage + jsonMessage)))

	messages,err := reader.Read()

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(messages) != 2 {
		t.Errorf("error parse message")
	}

	if messages[0].Action != "register" {
		t.Errorf("error parse action")
	}
}

func TestReadHalf(t *testing.T) {
	reader := NewReader(bytes.NewBuffer([]byte(jsonMessage + "hello world")))

	messages,err := reader.Read()

	if err != nil {
		t.Errorf(err.Error())
	}

	if messages[0].Action != "register" {
		t.Errorf("error parse action")
	}

	if reader.Incoming.String() != "hello world" {
		t.Errorf("error parse")
	}
}