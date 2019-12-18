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

// New returns a new instance of Repository
func New(db *ms.MemoryStore) Repository {
	return Repository{DB: db}
}

// CreateAccount a new event into the EventStore
func (ar *Repository) CreateAccount(ac Account) es.Event {
	time := time.Now()
	event := es.Event{Timestamp: time, Name: AccountCreated, Payload: ac}

	ar.DB.Append(event)
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
