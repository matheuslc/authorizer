package account

import (
	"sync"
	"testing"

	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestCreateAccount(t *testing.T) {
	account := Account{ActiveCard: true, AvailableLimit: 100}
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	repo := New(&memoryStore)

	repo.CreateAccount(account)

	events := memoryStore.Get()

	if events[0].Payload.(Account).ActiveCard == false {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestCreateConcurrent(t *testing.T) {
	times := 99
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	repo := New(&memoryStore)
	var wg sync.WaitGroup

	wg.Add(times)
	for i := 0; i < times; i++ {
		go createAndMarkAsDone(&repo, &wg)
	}
	wg.Wait()

	items := memoryStore.Get()

	if len(items) != 1 {
		t.Errorf("Concurrent try to insert a new account.")
	}
}

func createAndMarkAsDone(ar *Repository, wg *sync.WaitGroup) {
	account := Account{ActiveCard: true, AvailableLimit: 100}

	ar.CreateAccount(account)
	wg.Done()
}
