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

	transactionRepository := New(&memoryStore)
	transactionRepository.Append(transaction)

	events := memoryStore.Get()

	if events[0].Payload.(Transaction).ID != uuid {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestAppendConcurrent(t *testing.T) {
	namespace := "fake_name"
	memoryStore := ms.NewStorage(namespace)
	transactionRepository := New(&memoryStore)

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go appendAndMarkAsDone(&transactionRepository, &wg)
	}
	wg.Wait()

	items := memoryStore.Get()

	if len(items) != 100 {
		t.Errorf("Concurrent try to insert a new transaction.")
	}
}

func appendAndMarkAsDone(tr *TransactionRepository, wg *sync.WaitGroup) {
	uuid, _ := uuid.NewUUID()
	transaction := Transaction{ID: uuid, Merchant: "Carmels Store", Amount: 10, time: time.Now()}

	tr.Append(transaction)
	wg.Done()
}
