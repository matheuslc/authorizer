package account

import (
	"sync"
	"testing"

	es "github.com/matheuslc/authorizer/internal/eventstore"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestCreateAccount(t *testing.T) {
	account := Account{ActiveCard: true, AvailableLimit: 100}
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)
	repo := Repository{&memoryStore}

	accountEvent := es.Event{Payload: account}
	repo.CreateAccount(accountEvent)

	events := memoryStore.Get()
	if events[0].Payload.(Account).ActiveCard == false {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestCreateConcurrent(t *testing.T) {
	times := 99
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	repo := Repository{DB: &memoryStore}
	var wg sync.WaitGroup

	wg.Add(times)
	for i := 0; i < times; i++ {
		go createAndMarkAsDone(&repo, &wg)
	}
	wg.Wait()

	items := memoryStore.Get()

	if len(items) != times {
		t.Errorf("Can't insert all the items. %v was inserted.", len(items))
	}
}

func createAndMarkAsDone(ar *Repository, wg *sync.WaitGroup) {
	account := Account{ActiveCard: true, AvailableLimit: 100}
	event := es.Event{Payload: account}

	ar.CreateAccount(event)
	wg.Done()
}
