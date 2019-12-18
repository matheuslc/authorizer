package transactionentity

import (
	"testing"
	"time"

	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestAuthorizeTransactionSuccess(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 2000}

	trDone := es.Event{Payload: Transaction{Merchant: "nuevo-store", Amount: 1000}}
	trSuccess := Transaction{Merchant: "anonther-store", Amount: 700}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(trDone)

	useCase := AuthorizeUseCase{
		AccountRepo:       accountRepo,
		TransactionRepo:   transactionRepo,
		TransactionIntent: trSuccess,
	}

	event := useCase.Execute()
	events := transactionRepo.All()

	if len(event.Violations) > 0 {
		t.Errorf("Expected no violations to be throwed. Instead, %v was throwed", len(event.Violations))
	}

	if len(events) < 2 {
		t.Errorf("Expected 2 transactions in the event store")
	}
}

func TestAuthorizeViolation(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 200}

	trDone := Transaction{Merchant: "nuevo-store", Amount: 150}
	trWithoutLimit := Transaction{Merchant: "anonther-store", Amount: 70}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(es.Event{Name: TransactionValidated, Payload: trDone, Timestamp: time.Now()})

	useCase := AuthorizeUseCase{
		AccountRepo:       accountRepo,
		TransactionRepo:   transactionRepo,
		TransactionIntent: trWithoutLimit,
	}

	event := useCase.Execute()

	if len(event.Violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(event.Violations))
	}

	if event.Violations[0] != InsuficientLimit {
		t.Errorf("Expected %v as violation. Get %v", InsuficientLimit, event.Violations[0])
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
	transactionRepo.Append(es.Event{Payload: trDone, Timestamp: time.Now()})

	useCase := AuthorizeUseCase{
		AccountRepo:       accountRepo,
		TransactionRepo:   transactionRepo,
		TransactionIntent: trDoubled,
	}

	event := useCase.Execute()
	if len(event.Violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(event.Violations))
	}

	if event.Violations[0] != DoubledNotAllowed {
		t.Errorf("Expected %v as violation. Get %v", DoubledNotAllowed, event.Violations[0])
	}
}

func TestMoreThanAllowedTransactionViolationUseCase(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")

	accountRepo := ac.Repository{DB: &acStore}
	transactionRepo := Repository{DB: &tStore}

	account := ac.Account{ActiveCard: true, AvailableLimit: 1000}

	trNuevo := es.Event{Payload: Transaction{Merchant: "nuevo-store", Amount: 150}, Timestamp: time.Now()}
	trCarmel := es.Event{Payload: Transaction{Merchant: "carmel-store", Amount: 150}, Timestamp: time.Now()}
	trNubank := es.Event{Payload: Transaction{Merchant: "nubank-store", Amount: 150}, Timestamp: time.Now()}
	trExceded := Transaction{Merchant: "intent-store", Amount: 150}

	accountRepo.CreateAccount(account)
	transactionRepo.Append(trNuevo)
	transactionRepo.Append(trCarmel)
	transactionRepo.Append(trNubank)

	useCase := AuthorizeUseCase{
		AccountRepo:       accountRepo,
		TransactionRepo:   transactionRepo,
		TransactionIntent: trExceded,
	}

	event := useCase.Execute()

	if len(event.Violations) != 1 {
		t.Errorf("Expected just one violations and get %v", len(event.Violations))
	}

	if event.Violations[0] != MoreThanAllowed {
		t.Errorf("Expected %v as violation. Get %v", MoreThanAllowed, event.Violations[0])
	}
}
