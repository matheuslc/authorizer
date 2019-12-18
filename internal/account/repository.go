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

// RepositoryInterface defines the contract of an Account Repository
type RepositoryInterface interface {
	CreateAccount(ac Account) bool
	CurrentAccount() Account
}

// CreateAccount saves the new user inside the event store
func (ar *Repository) CreateAccount(event es.Event) {
	ar.DB.Append(event)
}

// NewEvent creates a new event
func (ar *Repository) NewEvent(ac Account, eventType string, violations []string) es.Event {
	time := time.Now()
	event := es.Event{Timestamp: time, Name: eventType, Payload: ac, Violations: violations}
	return event
}

// CurrentAccount returns the current system account
// Currently, we handle just one account. We loop through "account created" events searching for the account definition
func (ar *Repository) CurrentAccount() Account {
	c := ar.DB.EventsByName(AccountCreated)
	account := Account{}

	for event := range c {
		account = event.Payload.(Account)
	}

	return account
}

// Iter exposes all account events persisted inside the account store
func (ar *Repository) Iter() []es.Event {
	data := ar.DB.Iter()
	ocur := []es.Event{}

	for event := range data {
		ocur = append(ocur, event)
	}

	return ocur
}
