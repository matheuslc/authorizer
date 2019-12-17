package account

import (
	"testing"

	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestCreateUseCaseSuccess(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	accountRepo := Repository{DB: &acStore}

	account := Account{ActiveCard: true, AvailableLimit: 200}
	useCase := CreateUseCase{
		accountRepo:   accountRepo,
		accountIntent: account,
	}

	violations := useCase.Execute()

	if len(violations) > 0 {
		t.Errorf("Expected to not thrown any violation. Throwns: %v", len(violations))
	}
}

func TestCreateUseCaseViolation(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	accountRepo := Repository{DB: &acStore}

	account := Account{ActiveCard: true, AvailableLimit: 200}
	accountRepo.CreateAccount(account)

	useCase := CreateUseCase{
		accountRepo:   accountRepo,
		accountIntent: account,
	}

	violations := useCase.Execute()

	if len(violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(violations))
	}

	if violations[0] != AlreadyInitialized {
		t.Errorf("Expected %v as violation. Get %v", AlreadyInitialized, violations[0])
	}
}
