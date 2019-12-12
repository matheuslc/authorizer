package account

import (
	"time"

	"github.com/google/uuid"
	"github.com/matheuslc/authorizer/internal/eventstore"
)

// AccountRepository handle all writes on the account event store
type AccountRepository struct {
	db *eventstore.InMemoryStorage
}

// New returns a new instance of AccountRepository
func New(db *eventstore.InMemoryStorage) AccountRepository {
	return AccountRepository{db: db}
}

// CreateAccount a new event into the EventStore
func (ar *AccountRepository) CreateAccount(ac Account) bool {
	accountExists := ar.AccountAlreadyExists()
	if accountExists {
		return false
	}

	uuid, _ := uuid.NewUUID()
	time := time.Now()
	event := eventstore.Event{ID: uuid, Timestamp: time, Name: AccountCreated, Payload: ac}

	ar.db.Append(event)
	return true
}

// AccountAlreadyExists checks if an account already was initialized
func (ar *AccountRepository) AccountAlreadyExists() bool {
	events, _ := ar.db.Get()

	for _, event := range events {
		if event.Name == AccountCreated {
			return true
		}
	}

	return false
}
