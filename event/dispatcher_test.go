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