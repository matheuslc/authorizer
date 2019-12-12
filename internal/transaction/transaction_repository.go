package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/matheuslc/authorizer/internal/eventstore"
)

// TransactionRepository handle all writes on the account event store
type TransactionRepository struct {
	db *eventstore.InMemoryStorage
}

// New returns a new instance of TransactionRepository
func New(db *eventstore.InMemoryStorage) TransactionRepository {
	return TransactionRepository{db: db}
}

// Append a new event into the EventStore
func (tr *TransactionRepository) Append(t Transaction) bool {
	uuid, _ := uuid.NewUUID()
	time := time.Now()
	event := eventstore.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: t}

	tr.db.Append(event)
	return true
}
