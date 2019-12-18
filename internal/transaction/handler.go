package transaction

import (
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	"github.com/matheuslc/authorizer/internal/eventstore"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

// TransactionCommandHandler
func CommandHandler(
	payload map[string]interface{},
	acStore *memorystore.MemoryStore,
	trStore *memorystore.MemoryStore,
) eventstore.Event {
	acRepo := ac.Repository{DB: acStore}
	trRepo := Repository{DB: trStore}

	date, _ := time.Parse(time.RFC3339, payload["time"].(string))
	authUseCase := AuthorizeUseCase{
		AccountRepo:     acRepo,
		TransactionRepo: trRepo,
		TransactionIntent: Transaction{
			Merchant: payload["merchant"].(string),
			Amount:   int(payload["amount"].(float64)),
			Time:     date,
		},
	}

	availabeLimitUseCase := BalanceUseCase{TransactionRepo: trRepo}
	auth := authUseCase.Execute()
	total := availabeLimitUseCase.Execute()

	acc := acRepo.CurrentAccount()
	acc.AvailableLimit = acc.AvailableLimit - total

	return acRepo.NewEvent(acc, ac.AccountTransaction, auth.Violations)
}
