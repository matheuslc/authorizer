package transactionentity

import (
	"testing"

	ac "github.com/matheuslc/authorizer/internal/account"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestAuthorizeViolation(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 200}

	trDone := Transaction{Merchant: "nuevo-store", Amount: 150}
	trWithoutLimit := Transaction{Merchant: "anonther-store", Amount: 70}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(trDone)

	useCase := AuthorizeUseCase{
		accountRepo:       accountRepo,
		transactionRepo:   transactionRepo,
		transactionIntent: trWithoutLimit,
	}

	violations := useCase.Execute()

	if len(violations) == 2 {
		t.Errorf("lorem")
	}
}
