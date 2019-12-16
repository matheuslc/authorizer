package account

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

// New returns a new instance of Repository
func New(db *ms.MemoryStore) Repository {
	return Repository{DB: db}
}

// CreateAccount a new event into the EventStore
func (ar *Repository) CreateAccount(ac Account) bool {
	accountExists := ar.AccountAlreadyExists()
	if accountExists {
		return false
	}

	uuid, _ := uuid.NewUUID()
	time := time.Now()
	event := es.Event{ID: uuid, Timestamp: time, Name: AccountCreated, Payload: ac}

	ar.DB.Append(event)
	return true
}

// CurrentAccount
func (ar *Repository) CurrentAccount() Account {
	c := ar.DB.EventsByName("account:created")
	account := Account{}

	for event := range c {
		account = event.Payload.(Account)
	}

	return account
}

// AccountAlreadyExists checks if an account already was initialized
func (ar *Repository) AccountAlreadyExists() bool {
	c := ar.DB.Iter()

	for event := range c {
		if event.Name == AccountCreated {
			return true
		}
	}

	return false
}
