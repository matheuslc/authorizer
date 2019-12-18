package eventstore

import (
	"time"
)

// Event defines an event wrapper. Its payload contains the specific event
type Event struct {
	Timestamp  time.Time `json:"-"`
	Name       string    `json:"-"`
	Payload    interface{}
	Violations []string `json:"violations"`
}

// EventStore interface defines what an EventStore needs
type EventStore interface {
	Append(event Event) Event
	EventsByName(name string) <-chan Event
	Get() []Event
	Iter() <-chan Event
	IterAfter(after time.Time) <-chan Event
	NewStorage(namespace string) EventStore
}
