package account

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestCreateAccount(t *testing.T) {
	uuid, _ := uuid.NewUUID()

	account := Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	Repository := New(&memoryStore)

	Repository.CreateAccount(account)

	events := memoryStore.Get()

	if events[0].Payload.(Account).ActiveCard == false {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestCreateConcurrent(t *testing.T) {
	times := 99
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	Repository := New(&memoryStore)
	var wg sync.WaitGroup

	wg.Add(times)
	for i := 0; i < times; i++ {
		go createAndMarkAsDone(&Repository, &wg)
	}
	wg.Wait()

	items := memoryStore.Get()

	if len(items) != 1 {
		t.Errorf("Concurrent try to insert a new account.")
	}
}

func createAndMarkAsDone(ar *Repository, wg *sync.WaitGroup) {
	uuid, _ := uuid.NewUUID()
	account := Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	ar.CreateAccount(account)
	wg.Done()
}