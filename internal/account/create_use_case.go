package account

import es "github.com/matheuslc/authorizer/internal/eventstore"

// CreateUseCase struct defines what it takes to create an account
type CreateUseCase struct {
	AccountRepo   Repository
	AccountIntent Account
}

// Execute runs the violation functions under the account events and if no violations were found,
// then a new account is created
func (uc *CreateUseCase) Execute() es.Event {
	accountEvents := uc.AccountRepo.Iter()
	acViolation := Violations{AccountEvents: accountEvents, AccountIntent: uc.AccountIntent}

	violations := uc.runViolations(acViolation)
	event := uc.AccountRepo.NewEvent(uc.AccountIntent, AccountCreated, violations)

	if len(violations) == 0 {
		uc.AccountRepo.CreateAccount(event)
	}

	return event
}

func (uc *CreateUseCase) runViolations(acViolation Violations) []string {
	violations := []string{}

	already, alreadyStatus := AlreadyInitializedViolation(acViolation)
	active, activeStatus := ActiveCardViolation(acViolation)

	if alreadyStatus {
		violations = append(violations, already)
	}

	if activeStatus {
		violations = append(violations, active)
	}

	return violations
}
