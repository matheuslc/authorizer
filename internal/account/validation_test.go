package account

import (
	"testing"

	es "github.com/matheuslc/authorizer/internal/eventstore"
)

func TestAccountNotInitialized(t *testing.T) {
	ac := Account{}
	events := []es.Event{}

	tv := Violations{AccountIntent: ac, AccountEvents: events}
	_, ok := NotInitilizedViolation(tv)

	if !ok {
		t.Error("Violation failed, empty account was not catched")
	}
}

func TestAccountActiveCardViolation(t *testing.T) {
	ac := Account{ActiveCard: false, AvailableLimit: 100}
	events := []es.Event{}

	tv := Violations{AccountIntent: ac, AccountEvents: events}
	_, ok := ActiveCardViolation(tv)

	if !ok {
		t.Error("Violation failed, non ActiveCard was not catched")
	}
}

func TestAlreadyInitializedViolation(t *testing.T) {
	ac := Account{ActiveCard: false, AvailableLimit: 100}
	accountCreatedEvent := es.Event{Name: AccountCreated, Payload: ac}
	events := []es.Event{}

	tv := Violations{AccountIntent: ac, AccountEvents: append(events, accountCreatedEvent)}
	_, ok := AlreadyInitializedViolation(tv)

	if !ok {
		t.Error("Violation failed, non ActiveCard was not catched")
	}
}
