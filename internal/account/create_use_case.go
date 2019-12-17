package account

type CreateUseCase struct {
	accountRepo   Repository
	accountIntent Account
}

// Execute
func (uc *CreateUseCase) Execute() []string {
	accountEvents := uc.accountRepo.Iter()

	acViolation := Violations{AccountEvents: accountEvents, AccountIntent: uc.accountIntent}

	violations := uc.runViolations(acViolation)
	return violations
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
