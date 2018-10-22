package event

// Generic event struct
type Event struct {
	// Event name
	Name string

	// whether stop propagation
	PropagationStopped bool

	// parameters
	Parameters map[string]interface{}
}

// Create a new event
func NewEvent(name string, parameters map[string]interface{}) *Event{
	return &Event{
		Name: name,
		PropagationStopped: false,
		Parameters: parameters,
	}
}
