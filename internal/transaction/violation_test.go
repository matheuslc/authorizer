package transaction

import (
	"testing"
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

func TestMoreThanAllowedViolation(t *testing.T) {
	events := generateEvents(4)

	ac := ac.Account{ActiveCard: true, AvailableLimit: 100}
	time := time.Now()
	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 100, Time: time}

	tv := Violations{Account: ac, TransactionEvents: events, TransactionIntent: transaction}
	_, violations := MoreThanAllowedViolation(tv)

	if !violations {
		t.Error("Violation failed, expected: true. got: false")
	}
}

func TestMoreThanAllowedViolationSuccess(t *testing.T) {
	events := generateEvents(1)

	ac := ac.Account{ActiveCard: true, AvailableLimit: 100}
	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 100, Time: time.Now()}

	tv := Violations{Account: ac, TransactionEvents: events, TransactionIntent: transaction}
	_, violations := MoreThanAllowedViolation(tv)

	if violations {
		t.Error("Violation failerd, expected false, got true")
	}
}

func TestDuplicatedTransaction(t *testing.T) {
	events := generateEvents(1)

	ac := ac.Account{ActiveCard: true, AvailableLimit: 100}
	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 100, Time: time.Now()}

	tv := Violations{Account: ac, TransactionEvents: events, TransactionIntent: transaction}
	_, ok := DuplicatedTransaction(tv)

	if !ok {
		t.Error("Violation failed, expected to detect duplicae transaction.")
	}
}

func TestConcurrent(t *testing.T) {
	events := generateEvents(10)
	secondEvents := generateEvents(10)

	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 100, Time: time.Now()}
	ac := ac.Account{ActiveCard: true, AvailableLimit: 100}

	tv := Violations{Account: ac, TransactionEvents: events, TransactionIntent: transaction}
	tv2 := Violations{Account: ac, TransactionEvents: secondEvents, TransactionIntent: transaction}

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

	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 250, Time: time.Now()}
	ac := ac.Account{ActiveCard: true, AvailableLimit: 200}
	tv := Violations{Account: ac, TransactionEvents: events, TransactionIntent: transaction}

	_, ok := AccountLimitViolation(tv)

	if !ok {
		t.Error("Violation failed, account without limit was not catched")
	}
}

func genEvent() es.Event {
	time := time.Now()
	transaction := Transaction{Merchant: "carmel-restaurant", Amount: 100, Time: time}
	event := es.Event{Timestamp: time, Name: TransactionValidated, Payload: transaction}

	return event
}

func generateEvents(amount int) []es.Event {
	events := []es.Event{}

	for index := 0; index < amount; index++ {
		events = append(events, genEvent())
	}

	return events
}
