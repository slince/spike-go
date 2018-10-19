package event

// Listener type
type listener func(event *Event)

type ListenerQueue struct {
	listeners []listener
}

// Add a listener to queue.
func (q ListenerQueue) add(callback listener){
	q.listeners = append(q.listeners, callback)
}

// Add a listener to queue.
func (q ListenerQueue) remove(callback listener){
	for k, v := range q.listeners {
		if v == callback {

		}
	}
}


// Dispatcher
type Dispatcher struct {
	listeners map[string][]listener
}

// Add a listener to dispatcher.
func (dispatcher *Dispatcher) On(eventName string, callback listener) {

	listeners, ok := dispatcher.listeners[eventName]
	if !ok {
		listeners = make([]listener, 0, 10)
	}
	listeners = append(listeners, callback)
	dispatcher.listeners[eventName] = listeners
}

// Remove a listener from the dispatcher.
func (dispatcher *Dispatcher) Off(eventName string, callback listener) {
	if callback == nil {
		delete(dispatcher.listeners, eventName)
	} else if listeners, ok := dispatcher.listeners[eventName]; ok {
		
	}
}

// Fire the event
func (dispatcher *Dispatcher) Fire(event Event){

	listeners, ok := dispatcher.listeners[event.GetName()]

	if ok {
		for _, listener := range listeners {
			listener(&event)
		}
	}
}

// Creates a new event dispatcher.
func NewDispatcher() *Dispatcher{
	return &Dispatcher{
		make(map[string][]listener),
	}
}

