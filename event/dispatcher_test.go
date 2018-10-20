package event

import (
	"testing"
)


func TestOn(t *testing.T) {

	var dispatcher = NewDispatcher()
	var foo  = 10

	var callback = func(event *Event) {
		foo ++
	}

	dispatcher.On("foo", (*Listener)(&callback))
	
	var event = NewEvent("foo", nil)
	dispatcher.Fire(event)

	if foo != 11 {
		t.Errorf("The Listener is not fired")
	}
}



func TestOff(t *testing.T) {

	var dispatcher = NewDispatcher()
	var foo  = 10

	//dispatcher.On("foo", func1)
	var lis = NewListener(func (event *Event) {
		foo ++
	})

	dispatcher.On("foo", lis)
	dispatcher.Off("foo", lis)

	var event = NewEvent("foo", nil)
	dispatcher.Fire(event)

	if foo == 11 {
		t.Errorf("Off error")
	}
}

func TestAddSubscriber(t *testing.T) {

	var dispatcher = NewDispatcher()
	var foo,bar  = 10, 11

	var sub = NewSubscriber(map[string]func(event *Event){
		"foo": func(event *Event) {
			foo ++
		},
		"bar": func(event *Event) {
			bar ++
		},
	})

	dispatcher.AddSubscriber(sub)
	dispatcher.Fire(NewEvent("foo", nil))
	dispatcher.Fire(NewEvent("bar", nil))

	if foo != 11 || bar != 12 {
		t.Errorf("Off error")
	}
}

func TestRemoveSubscriber(t *testing.T) {

	var dispatcher = NewDispatcher()
	var foo,bar  = 10, 11

	var sub = NewSubscriber(map[string]func(event *Event){
		"foo": func(event *Event) {
			foo ++
		},
		"bar": func(event *Event) {
			bar ++
		},
	})

	dispatcher.AddSubscriber(sub)
	dispatcher.RemoveSubscriber(sub)

	dispatcher.Fire(NewEvent("foo", nil))
	dispatcher.Fire(NewEvent("bar", nil))

	if foo == 11 || bar == 12 {
		t.Errorf("Off error")
	}
}