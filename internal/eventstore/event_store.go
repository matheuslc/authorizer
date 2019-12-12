package eventstore

import (
	"time"

	"github.com/google/uuid"
)

// Event defines an event wrapper. Its payload contains the specific event
type Event struct {
	ID        uuid.UUID
	Timestamp time.Time
	Name      string
	Payload   interface{}
}

// EventStore interface defines what an EventStore needs
type EventStore interface {
	Append(event Event)
	Get() []Event
	Iter() <-chan Event
	NewStorage(namespace string) EventStore
}
