package transactionentity

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	ms "github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TesteAppendTransaction(t *testing.T) {
	uuid, _ := uuid.NewUUID()

	transaction := Transaction{ID: uuid, Merchant: "Carmels Store", Amount: 10, time: time.Now()}
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)

	Repository := New(&memoryStore)
	Repository.Append(transaction, []string{})

	events := memoryStore.Get()

	if events[0].Payload.(Transaction).ID != uuid {
		t.Errorf("Event out of order was inserted sorted")
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
	uuid, _ := uuid.NewUUID()
	transaction := Transaction{ID: uuid, Merchant: "Carmels Store", Amount: 10, time: time.Now()}

	tr.Append(transaction, []string{})
	wg.Done()
}
