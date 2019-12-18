package transactionentity

import (
	"sync"
	"testing"
	"time"

	es "github.com/matheuslc/authorizer/internal/eventstore"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TesteAppendTransaction(t *testing.T) {
	transaction := Transaction{Merchant: "Carmels Store", Amount: 10, Time: time.Now()}
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	repo := New(&memoryStore)
	repo.Append(es.Event{Payload: transaction})

	events := memoryStore.Get()

	if events[0].Payload.(Transaction).Amount != 10 {
		t.Errorf("Event was not inserted as expected")
	}
}

func TestAppendConcurrent(t *testing.T) {
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)
	Repository := New(&memoryStore)

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go appendAndMarkAsDone(&Repository, &wg)
	}
	wg.Wait()

	items := memoryStore.Get()

	if len(items) != 100 {
		t.Errorf("Concurrent try to insert a new transaction.")
	}
}

func appendAndMarkAsDone(tr *Repository, wg *sync.WaitGroup) {
	transaction := Transaction{Merchant: "Carmels Store", Amount: 10, Time: time.Now()}

	tr.Append(es.Event{Payload: transaction})
	wg.Done()
}
