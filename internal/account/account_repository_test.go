package account

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/matheuslc/authorizer/internal/eventstore"
)

func TestCreateAccount(t *testing.T) {
	uuid, _ := uuid.NewUUID()

	account := Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}
	namespace := "fake_name"
	inMemoryStorage := eventstore.NewStorage(namespace)

	accountRepository := New(&inMemoryStorage)

	accountRepository.CreateAccount(account)

	events, ok := inMemoryStorage.Get()

	if ok == false {
		t.Errorf("Event out of order was inserted sorted")
	}

	if events[0].Payload.(Account).ActiveCard == false {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestCreateConcurrent(t *testing.T) {

	namespace := "fake_name"
	inMemoryStorage := eventstore.NewStorage(namespace)

	accountRepository := New(&inMemoryStorage)
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go createAndMarkAsDone(&accountRepository, &wg)
	}
	wg.Wait()

	items, _ := inMemoryStorage.Get()

	if len(items) != 1 {
		t.Errorf("Concurrent try to insert a new account.")
	}
}

func createAndMarkAsDone(ar *AccountRepository, wg *sync.WaitGroup) {
	uuid, _ := uuid.NewUUID()
	account := Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	ar.CreateAccount(account)
	wg.Done()
}
