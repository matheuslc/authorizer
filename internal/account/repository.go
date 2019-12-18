package account

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
	CreateAccount(ac Account) bool
	CurrentAccount() Account
}

// CreateAccount a new event into the EventStore
func (ar *Repository) CreateAccount(event es.Event) {
	ar.DB.Append(event)
}

// NewEvent
func (ar *Repository) NewEvent(ac Account, eventType string, violations []string) es.Event {
	time := time.Now()
	event := es.Event{Timestamp: time, Name: eventType, Payload: ac, Violations: violations}
	return event
}

// CurrentAccount
func (ar *Repository) CurrentAccount() Account {
	c := ar.DB.EventsByName(AccountCreated)
	account := Account{}

	for event := range c {
		account = event.Payload.(Account)
	}

	return account
}

// Iter
func (ar *Repository) Iter() []es.Event {
	data := ar.DB.Iter()
	ocur := []es.Event{}

	for event := range data {
		ocur = append(ocur, event)
	}

	return ocur
}
