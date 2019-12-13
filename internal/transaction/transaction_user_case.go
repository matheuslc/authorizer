package transactionentity

import (
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// AuthorizeTransactionUseCase
type AuthorizeTransactionUseCase struct {
	ur ac.AccountRepository
	tr TransactionRepository
	t  Transaction
}

// Execute
func (uc *AuthorizeTransactionUseCase) Execute() []string {
	now := time.Now()
	twoMinutes := time.Minute + time.Duration(2)
	past := now.Add(-twoMinutes)

	account := ac.Account{ActiveCard: true, AvailableLimit: 200}
	events := uc.tr.IterAfter(past)

	event := es.Event{Name: TransactionValidated, Payload: uc.t}
	tv := TransactionValidation{User: account, TransactionEvents: events, CurrentEvent: event}

	violations := []string{}

	allowed, allowedStatus := MoreThanAllowedViolation(tv)
	duplicated, duplicatedStatus := DuplicatedTransaction(tv)
	limit, limitStatus := AccountLimitViolation(tv)
	init, initStatus := AccountNotInitilizedViolation(tv)
	active, activeStatus := AccountActiveCardViolation(tv)

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
