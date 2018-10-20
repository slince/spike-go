package event

import (
	"testing"
)

var dispatcher = NewDispatcher()

func TestOn(t *testing.T) {

	var foo  = 10

	dispatcher.On("foo", func(event *Event) {
		foo ++
	})
	
	var event = NewEvent("foo", nil)
	dispatcher.Fire(event)

	if foo != 11 {
		t.Errorf("The listener is not fired")
	}
}

func TestOff(t *testing.T) {
	var foo  = 10

	var func1 = func(event *Event) {
		foo ++
	}

	dispatcher.On("foo", func1)
	dispatcher.Off("foo", &func1)

	var event = NewEvent("foo", nil)
	dispatcher.Fire(event)

	if foo == 11 {
		t.Errorf("Off error")
	}
}