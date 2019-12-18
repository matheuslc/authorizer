package transactionentity

import (
	"time"

	es "github.com/matheuslc/authorizer/internal/eventstore"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

// Repository handle all writes on the account event store
type Repository struct {
	DB *ms.MemoryStore
}

// RepositoryInterface
type RepositoryInterface interface {
	Append(t Transaction, v []string) es.Event
	IterAfter(after time.Time) []es.Event
}

// New returns a new instance of Repository
func New(db *ms.MemoryStore) Repository {
	return Repository{DB: db}
}

// NewEvent
func (tr *Repository) NewEvent(t Transaction, evenType string, violations []string) es.Event {
	time := time.Now()

	event := es.Event{
		Timestamp:  time,
		Name:       evenType,
		Payload:    t,
		Violations: violations,
	}

	return event
}

// Append a new event into the EventStore
func (tr *Repository) Append(ev es.Event) {
	tr.DB.Append(ev)
}

// All
func (tr *Repository) All() []es.Event {
	data := tr.DB.Get()
	return data
}

// IterAfter iterates thogh events after a certain time
func (tr *Repository) IterAfter(after time.Time) []es.Event {
	data := tr.DB.IterAfter(after)
	ocur := []es.Event{}

	for event := range data {
		if event.Timestamp.After(after) || event.Timestamp.Equal(after) {
			ocur = append(ocur, event)
		}
	}

	return ocur
}
