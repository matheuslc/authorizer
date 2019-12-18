package transaction

// BalanceUseCase
type BalanceUseCase struct {
	TransactionRepo Repository
}

// Execute
func (uc *BalanceUseCase) Execute() int {
	totalSpent := 0
	transctionEvents := uc.TransactionRepo.All()

	for _, event := range transctionEvents {

		if len(event.Violations) == 0 {
			t := event.Payload.(Transaction)
			totalSpent += t.Amount
		}
	}

	return totalSpent
}
