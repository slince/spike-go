package event

type Event interface {
	// Get the event name
	GetName() string

	// Stop propagation
	StopPropagation()

	// Checks whether propagation was stopped.
	IsPropagationStopped() bool
}

// Generic event struct
type GenericEvent struct {
	name string

	propagationStopped bool

	Data map[string]interface{}
}

// Get the event name
func (event *GenericEvent) GetName() string {
	return event.name
}

// Stop propagation
func (event *GenericEvent) StopPropagation() {
	event.propagationStopped = true
}

// Checks whether propagation was stopped.
func (event *GenericEvent) IsPropagationStopped() bool {
	return event.propagationStopped
}

// Create a new event
func NewEvent(name string, data map[string]interface{}) Event{
	return &GenericEvent{
		name: name,
		propagationStopped: false,
		Data: data,
	}
}
