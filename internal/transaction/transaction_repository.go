package transactionentity

import (
	"time"

	"github.com/google/uuid"
	es "github.com/matheuslc/authorizer/internal/eventstore"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

// TransactionRepository handle all writes on the account event store
type TransactionRepository struct {
	db *ms.MemoryStore
}

// New returns a new instance of TransactionRepository
func New(db *ms.MemoryStore) TransactionRepository {
	return TransactionRepository{db: db}
}

// Append a new event into the EventStore
func (tr *TransactionRepository) Append(t Transaction) bool {
	uuid, _ := uuid.NewUUID()
	time := time.Now()
	event := es.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: t}

	tr.db.Append(event)
	return true
}
