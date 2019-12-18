package account

import es "github.com/matheuslc/authorizer/internal/eventstore"

// CreateUseCase
type CreateUseCase struct {
	AccountRepo   Repository
	AccountIntent Account
}

// Execute
func (uc *CreateUseCase) Execute() es.Event {
	// accountEvents := uc.AccountRepo.Iter()

	// acViolation := Violations{AccountEvents: accountEvents, AccountIntent: uc.AccountIntent}

	// violations := uc.runViolations(acViolation)
	event := uc.AccountRepo.CreateAccount(uc.AccountIntent)

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
