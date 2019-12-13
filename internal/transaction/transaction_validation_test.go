package transactionentity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

func TestMoreThanAllowedViolation(t *testing.T) {
	events := generateEvents(4)
	uuid, _ := uuid.NewUUID()

	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}
	time := time.Now()
	transaction := Transaction{ID: uuid, Merchant: "carmel-restaurant", Amount: 100, time: time}
	event := es.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: transaction}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: event}
	violations := MoreThanAllowedViolation(tv)

	if !violations {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func TestMoreThanAllowedViolationSuccess(t *testing.T) {
	events := generateEvents(1)
	event := genEvent()
	uuid, _ := uuid.NewUUID()
	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: event}
	violations := MoreThanAllowedViolation(tv)

	if violations {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func TestDuplicatedTransaction(t *testing.T) {
	events := generateEvents(1)
	event := genEvent()
	uuid, _ := uuid.NewUUID()
	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: event}
	_, ok := DuplicatedTransaction(tv)

	if !ok {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func TestConcurrent(t *testing.T) {
	events := generateEvents(10)
	secondEvents := generateEvents(10)
	firstEvent := genEvent()
	uuid, _ := uuid.NewUUID()
	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: firstEvent}
	tv2 := TransactionValidation{User: user, TransactionEvents: secondEvents, CurrentEvent: firstEvent}

	value := MoreThanAllowedViolation(tv)
	_, ok := DuplicatedTransaction(tv2)

	if !value {
		t.Error("More than allowed violation failed with a bucket of four events")
	}

	if !ok {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func TestAccountLimitViolation(t *testing.T) {
	events := generateEvents(10)
	firstEvent := genEvent()
	uuid, _ := uuid.NewUUID()
	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 600}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: firstEvent}

	_, ok := AccountLimitViolation(tv)

	if !ok {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func TestAccountNotInitialized(t *testing.T) {
	events := generateEvents(2)
	firstEvent := genEvent()
	user := ac.Account{}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: firstEvent}

	_, ok := AccountNotInitilizedViolation(tv)

	if !ok {
		t.Error("More than allowed violation failed with a bucket of four events")
	}
}

func genEvent() es.Event {
	uuid, _ := uuid.NewUUID()
	time := time.Now()
	transaction := Transaction{ID: uuid, Merchant: "carmel-restaurant", Amount: 100, time: time}
	event := es.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: transaction}

	return event
}

func generateEvents(amount int) <-chan es.Event {
	c := make(chan es.Event)

	go func() {
		defer close(c)

		uuid, _ := uuid.NewUUID()
		time := time.Now()
		transaction := Transaction{ID: uuid, Merchant: "carmel-restaurant", Amount: 100, time: time}
		event := es.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: transaction}

		for i := 0; i < amount; i++ {
			c <- event
		}
	}()

	return c
}