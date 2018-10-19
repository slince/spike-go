package event

// Listener type
type listener func(event *Event)

type ListenerQueue struct {
	listeners []*listener
}

// Add a listener to queue.
func (q *ListenerQueue) add(callback *listener){
	q.listeners = append(q.listeners, callback)
}

// Add a listener to queue.
func (q *ListenerQueue) remove(callback *listener){
	for k, v := range q.listeners {
		if v == callback {
			q.listeners = append(q.listeners[:k], q.listeners[k+1:]...)
			return
		}
	}
}


// Dispatcher
type Dispatcher struct {
	listeners map[string]*ListenerQueue
}

// Add a listener to dispatcher.
func (dispatcher *Dispatcher) On(eventName string, callback listener) {
	listenerQueue, ok := dispatcher.listeners[eventName]
	if !ok {
		listenerQueue = &ListenerQueue{
			make([]*listener, 0),
		}
		dispatcher.listeners[eventName] = listenerQueue
	}
	listenerQueue.add(&callback)
}

// Remove a listener from the dispatcher.
func (dispatcher *Dispatcher) Off(eventName string, callback listener) {
	if callback == nil {
		delete(dispatcher.listeners, eventName)
	} else if listenerQueue, ok := dispatcher.listeners[eventName]; ok {
		listenerQueue.remove(&callback)
	}
}

// Fire the event
func (dispatcher *Dispatcher) Fire(event Event){
	listenerQueue, ok := dispatcher.listeners[event.GetName()]

	if ok {
		for _, callback := range listenerQueue.listeners {
			(*callback)(&event)
		}
	}
}

// Creates a new event dispatcher.
func NewDispatcher() *Dispatcher{
	return &Dispatcher{
		make(map[string]*ListenerQueue),
	}
}

