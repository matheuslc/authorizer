package transaction

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matheuslc/authorizer/internal/eventstore"
)

func TesteAppendTransaction(t *testing.T) {
	uuid, _ := uuid.NewUUID()

	transaction := Transaction{ID: uuid, Merchant: "Carmels Store", Amount: 10, time: time.Now()}
	namespace := "fake_name"
	MemoryStore := eventstore.NewStorage(namespace)

	transactionRepository := New(&MemoryStore)
	transactionRepository.Append(transaction)

	events, ok := MemoryStore.Get()

	if ok == false {
		t.Errorf("Event out of order was inserted sorted")
	}

	if events[0].Payload.(Transaction).ID != uuid {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestAppendConcurrent(t *testing.T) {
	namespace := "fake_name"
	MemoryStore := eventstore.NewStorage(namespace)
	transactionRepository := New(&MemoryStore)

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go appendAndMarkAsDone(&transactionRepository, &wg)
	}
	wg.Wait()

	items, _ := MemoryStore.Get()

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
