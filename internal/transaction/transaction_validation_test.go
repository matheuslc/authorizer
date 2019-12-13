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
	_, violations := MoreThanAllowedViolation(tv)

	if !violations {
		t.Error("Violation failed, expected: true. got: false")
	}
}

func TestMoreThanAllowedViolationSuccess(t *testing.T) {
	events := generateEvents(1)
	event := genEvent()
	uuid, _ := uuid.NewUUID()
	user := ac.Account{ID: uuid, ActiveCard: true, AvailableLimit: 100}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: event}
	_, violations := MoreThanAllowedViolation(tv)

	if violations {
		t.Error("Violation failerd, expected false, got true")
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
		t.Error("Violation failed, expected to detect duplicae transaction.")
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

	_, value := MoreThanAllowedViolation(tv)
	_, ok := DuplicatedTransaction(tv2)

	if !value {
		t.Error("MoreThanAllowedViolation failed running concurrent")
	}

	if !ok {
		t.Error("DuplicatedTransaction failed running concurrent")
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
		t.Error("Violation failed, user without limit was not catched")
	}
}

func TestAccountNotInitialized(t *testing.T) {
	events := generateEvents(2)
	firstEvent := genEvent()
	user := ac.Account{}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: firstEvent}
	_, ok := AccountNotInitilizedViolation(tv)

	if !ok {
		t.Error("Violation failerd, empty account was not catched")
	}
}

func TestAccountActiveCardViolation(t *testing.T) {
	events := generateEvents(2)
	firstEvent := genEvent()

	user := ac.Account{ActiveCard: false, AvailableLimit: 100}

	tv := TransactionValidation{User: user, TransactionEvents: events, CurrentEvent: firstEvent}
	_, ok := AccountActiveCardViolation(tv)

	if !ok {
		t.Error("Violation failed, non ActiveCard was not catched")
	}
}

func genEvent() es.Event {
	uuid, _ := uuid.NewUUID()
	time := time.Now()
	transaction := Transaction{ID: uuid, Merchant: "carmel-restaurant", Amount: 100, time: time}
	event := es.Event{ID: uuid, Timestamp: time, Name: TransactionValidated, Payload: transaction}

	return event
}

func generateEvents(amount int) []es.Event {
	events := []es.Event{}

	for index := 0; index < amount; index++ {
		events = append(events, genEvent())
	}

	return events
}
