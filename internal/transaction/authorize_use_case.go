package transactionentity

import (
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// AuthorizeUseCase
type AuthorizeUseCase struct {
	ur ac.Repository
	tr Repository
	t  Transaction
}

// Execute
func (uc *AuthorizeUseCase) Execute() []string {
	now := time.Now()
	twoMinutes := time.Minute + time.Duration(2)
	past := now.Add(-twoMinutes)

	account := ac.Account{ActiveCard: true, AvailableLimit: 200}
	events := uc.tr.IterAfter(past)

	event := es.Event{Name: TransactionValidated, Payload: uc.t}

	acViolation := ac.Violations{Account: account}
	trViolation := Violations{Account: account, TransactionEvents: events, CurrentEvent: event}

	violations := []string{}

	init, initStatus := ac.NotInitilizedViolation(acViolation)
	active, activeStatus := ac.ActiveCardViolation(acViolation)

	allowed, allowedStatus := MoreThanAllowedViolation(trViolation)
	duplicated, duplicatedStatus := DuplicatedTransaction(trViolation)
	limit, limitStatus := AccountLimitViolation(trViolation)

	if allowedStatus {
		violations = append(violations, allowed)
	}

	if duplicatedStatus {
		violations = append(violations, duplicated)
	}

	if limitStatus {
		violations = append(violations, limit)
	}

	if initStatus {
		violations = append(violations, init)
	}

	if activeStatus {
		violations = append(violations, active)
	}

	return violations
}
