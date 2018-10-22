package event

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent("init", nil)

	if _, ok := interface{}(event).(*Event); !ok {
		t.Errorf("create event error")
	}

	if event.Name != "init" {
		t.Errorf("bad get name")
	}
}


func TestEvent(t *testing.T) {
	event := NewEvent("init", nil)

	if event.PropagationStopped {
		t.Errorf("error init data")
	}

	event.PropagationStopped = true

	if !event.PropagationStopped {
		t.Errorf("error stop propagation")
	}

}
