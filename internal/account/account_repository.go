package account

import (
	"time"

	"github.com/google/uuid"
	es "github.com/matheuslc/authorizer/internal/eventstore"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

// AccountRepository handle all writes on the account event store
type AccountRepository struct {
	db *ms.MemoryStore
}

// New returns a new instance of AccountRepository
func New(db *ms.MemoryStore) AccountRepository {
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
	event := es.Event{ID: uuid, Timestamp: time, Name: AccountCreated, Payload: ac}

	ar.db.Append(event)
	return true
}

// CurrentAccount
func (ar *AccountRepository) CurrentAccount(ac Account) Account {
	c := ar.db.EventsByName("account:created")
	account := Account{}

	for event := range c {
		account = event.Payload.(Account)
	}

	return account
}

// AccountAlreadyExists checks if an account already was initialized
func (ar *AccountRepository) AccountAlreadyExists() bool {
	c := ar.db.Iter()

	for event := range c {
		if event.Name == AccountCreated {
			return true
		}
	}

	return false
}
