package event

// Listener type
type Listener func(ev *Event)

// Subscriber
type Subscriber struct {
	Listeners map[string]*Listener
}

// Listener queue
type listenerQueue struct {
	listeners []*Listener
}

// Add a Listener to queue.
func (lq *listenerQueue) add(callback *Listener){
	lq.listeners = append(lq.listeners, callback)
}

// Add a Listener to queue.
func (lq *listenerQueue) remove(callback *Listener){
	for k, v := range lq.listeners {
		if v == callback {
			lq.listeners = append(lq.listeners[:k], lq.listeners[k+1:]...)
			return
		}
	}
}

// Dispatcher
type Dispatcher struct {
	listeners map[string]*listenerQueue
}

// Add a Subscriber
func (dispatcher *Dispatcher) AddSubscriber(sub *Subscriber) {
	for ev, callback := range sub.Listeners {
		dispatcher.On(ev, callback)
	}
}

// Remove a Subscriber
func (dispatcher *Dispatcher) RemoveSubscriber(sub *Subscriber) {
	for ev, callback := range sub.Listeners {
		dispatcher.Off(ev, callback)
	}
}

// Add a Listener to dispatcher.
func (dispatcher *Dispatcher) On(eventName string, callback *Listener) {
	lq, ok := dispatcher.listeners[eventName]
	if !ok {
		lq = &listenerQueue{
			make([]*Listener, 0),
		}
		dispatcher.listeners[eventName] = lq
	}
	lq.add(callback)
}

// Remove a Listener from the dispatcher.
func (dispatcher *Dispatcher) Off(eventName string, callback *Listener) {
	if callback == nil {
		delete(dispatcher.listeners, eventName)
	} else if listenerQueue, ok := dispatcher.listeners[eventName]; ok {
		listenerQueue.remove(callback)
	}
}

// Fire the event
func (dispatcher *Dispatcher) Fire(event *Event){
	lq, ok := dispatcher.listeners[event.Name]

	if ok {
		for _, callback := range lq.listeners {
			(*callback)(event)
		}
	}
}

// Creates a new event dispatcher.
func NewDispatcher() *Dispatcher{
	return &Dispatcher{
		make(map[string]*listenerQueue),
	}
}

// Creates a new listener.
func NewListener(callback func(event *Event)) *Listener{
	return (*Listener)(&callback)
}


// Creates a new listener.
func NewSubscriber(listeners map[string]func(event *Event)) *Subscriber{

	var _listeners = make(map[string]*Listener, len(listeners))

	for eventName, callback := range listeners {
		_listeners[eventName] = NewListener(callback)
	}

	return &Subscriber{
		_listeners,
	}
}

