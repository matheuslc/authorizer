package account

import "testing"

func TestAccountNotInitialized(t *testing.T) {
	ac := Account{}

	tv := Violations{Account: ac}
	_, ok := NotInitilizedViolation(tv)

	if !ok {
		t.Error("Violation failerd, empty account was not catched")
	}
}

func TestAccountActiveCardViolation(t *testing.T) {
	ac := Account{ActiveCard: false, AvailableLimit: 100}

	tv := Violations{Account: ac}
	_, ok := ActiveCardViolation(tv)

	if !ok {
		t.Error("Violation failed, non ActiveCard was not catched")
	}
}
