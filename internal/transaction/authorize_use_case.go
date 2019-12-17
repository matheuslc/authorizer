package transactionentity

import (
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// AuthorizeUseCase
type AuthorizeUseCase struct {
	accountRepo       ac.Repository
	transactionRepo   Repository
	transactionIntent Transaction
}

// Execute
func (uc *AuthorizeUseCase) Execute() []string {
	account := uc.accountRepo.CurrentAccount()
	events := uc.transactionRepo.IterAfter(uc.pastDateToGetEvents())

	violations := uc.runViolations(account, events)
	return violations
}

func (uc *AuthorizeUseCase) runViolations(account ac.Account, events []es.Event) []string {
	violations := []string{}

	acViolation := ac.Violations{Account: account}
	trViolation := Violations{Account: account, TransactionEvents: events, TransactionIntent: uc.transactionIntent}

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
