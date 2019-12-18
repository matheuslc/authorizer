package transactionentity

import (
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// AuthorizeUseCase
type AuthorizeUseCase struct {
	AccountRepo       ac.Repository
	TransactionRepo   Repository
	TransactionIntent Transaction
}

// Execute
func (uc *AuthorizeUseCase) Execute() es.Event {
	account := uc.AccountRepo.CurrentAccount()
	transctionEvents := uc.TransactionRepo.IterAfter(uc.pastDateToGetEvents())
	accountEvents := uc.AccountRepo.Iter()

	acViolation := ac.Violations{AccountIntent: account, AccountEvents: accountEvents}
	trViolation := Violations{Account: account, TransactionEvents: transctionEvents, TransactionIntent: uc.TransactionIntent}

	violations := uc.runViolations(acViolation, trViolation)
	event := uc.TransactionRepo.NewEvent(uc.TransactionIntent, TransactionValidated, violations)

	uc.TransactionRepo.Append(event)
	return event
}

func (uc *AuthorizeUseCase) runViolations(acViolation ac.Violations, trViolation Violations) []string {
	violations := []string{}

	init, initStatus := ac.NotInitilizedViolation(acViolation)
	active, activeStatus := ac.ActiveCardViolation(acViolation)

	allowed, allowedStatus := MoreThanAllowedViolation(trViolation)
	duplicated, duplicatedStatus := DuplicatedTransaction(trViolation)
	limit, limitStatus := AccountLimitViolation(trViolation)

	if initStatus {
		violations = append(violations, init)
	}

	if activeStatus {
		violations = append(violations, active)
	}

	if allowedStatus {
		violations = append(violations, allowed)
	}

	if duplicatedStatus {
		violations = append(violations, duplicated)
	}

	if limitStatus {
		violations = append(violations, limit)
	}

	return violations
}

func (uc *AuthorizeUseCase) pastDateToGetEvents() time.Time {
	now := time.Now()
	twoMinutes := time.Minute + time.Duration(2)
	past := now.Add(-twoMinutes)

	return past
}
