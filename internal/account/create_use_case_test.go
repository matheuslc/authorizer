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
		AccountRepo:   accountRepo,
		AccountIntent: account,
	}

	event := useCase.Execute()

	if len(event.Violations) > 0 {
		t.Errorf("Expected to not thrown any violation. Throwns: %v", len(event.Violations))
	}
}

func TestCreateUseCaseViolation(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	accountRepo := Repository{DB: &acStore}

	account := Account{ActiveCard: true, AvailableLimit: 200}
	accountRepo.CreateAccount(account)

	useCase := CreateUseCase{
		AccountRepo:   accountRepo,
		AccountIntent: account,
	}

	event := useCase.Execute()

	if len(event.Violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(event.Violations))
	}

	if event.Violations[0] != AlreadyInitialized {
		t.Errorf("Expected %v as violation. Get %v", AlreadyInitialized, event.Violations[0])
	}
}
