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

	if len(violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(violations))
	}

	if violations[0] != InsuficientLimit {
		t.Errorf("Expected %v as violation. Get %v", InsuficientLimit, violations[0])
	}
}

func TestDoubleTransactionViolationUseCase(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 1000}

	trDone := Transaction{Merchant: "nuevo-store", Amount: 150}
	trDoubled := Transaction{Merchant: "nuevo-store", Amount: 150}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(trDone)

	useCase := AuthorizeUseCase{
		accountRepo:       accountRepo,
		transactionRepo:   transactionRepo,
		transactionIntent: trDoubled,
	}

	violations := useCase.Execute()

	if len(violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(violations))
	}

	if violations[0] != DoubledNotAllowed {
		t.Errorf("Expected %v as violation. Get %v", DoubledNotAllowed, violations[0])
	}
}

func TestMoreThanAllowedTransactionViolationUseCase(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 1000}

	trNuevo := Transaction{Merchant: "nuevo-store", Amount: 150}
	trCarmel := Transaction{Merchant: "carmel-store", Amount: 150}
	trNubank := Transaction{Merchant: "nubank-store", Amount: 150}
	trExceded := Transaction{Merchant: "intent-store", Amount: 150}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(trNuevo)
	transactionRepo.Append(trCarmel)
	transactionRepo.Append(trNubank)

	useCase := AuthorizeUseCase{
		accountRepo:       accountRepo,
		transactionRepo:   transactionRepo,
		transactionIntent: trExceded,
	}

	violations := useCase.Execute()

	if len(violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(violations))
	}

	if violations[0] != MoreThanAllowed {
		t.Errorf("Expected %v as violation. Get %v", MoreThanAllowed, violations[0])
	}
}
