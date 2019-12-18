package transactionentity

import (
	"time"

	"github.com/google/uuid"
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

// Append a new event into the EventStore
func (tr *Repository) Append(t Transaction, v []string) es.Event {
	uuid, _ := uuid.NewUUID()
	time := time.Now()

	event := es.Event{
		ID:         uuid,
		Timestamp:  time,
		Name:       TransactionValidated,
		Payload:    t,
		Violations: v,
	}

	tr.DB.Append(event)
	return event
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
