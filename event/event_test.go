package event

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent("init", nil)

	if _, ok := event.(Event); !ok {
		t.Errorf("create event error")
	}

	if event.GetName() != "init" {
		t.Errorf("bad get name")
	}
}


func TestEvent(t *testing.T) {
	event := NewEvent("init", nil)

	if event.IsPropagationStopped() {
		t.Errorf("error init data")
	}

	event.StopPropagation()

	if !event.IsPropagationStopped() {
		t.Errorf("error stop propagation")
	}

}
