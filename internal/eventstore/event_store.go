package eventstore

import (
	"time"
)

// Event defines an event wrapper. Its payload contains the specific event
type Event struct {
	Timestamp  time.Time
	Name       string
	Payload    interface{}
	Violations []string
}

// EventStore interface defines what an EventStore needs
type EventStore interface {
	Append(event Event)
	Get() []Event
	Iter() <-chan Event
	NewStorage(namespace string) EventStore
}

// EventJSONPresenter
type EventJSONPresenter struct {
	Event
	Violations []string
}
