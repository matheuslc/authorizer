package transactionentity

import (
	"testing"

	"github.com/matheuslc/authorizer/internal/account"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

func TestAuthorizeAccountAccountViolation(t *testing.T) {
	acStore := memorystore.NewStorage("account")
	tStore := memorystore.NewStorage("transaction")
	acRepository := account.AccountRepository{DB: &acStore}
	tRepository := TransactionRepository{DB: &tStore}

	tr := Transaction{Merchant: "nuevo-store", Amount: 200}
	tr2 := Transaction{Merchant: "anonther-store", Amount: 20}
	tr3 := Transaction{Merchant: "anonther-store", Amount: 22}
	tr4 := Transaction{Merchant: "anonther-store", Amount: 22}

	tRepository.Append(tr)
	tRepository.Append(tr2)
	tRepository.Append(tr3)
	tRepository.Append(tr4)

	useCase := AuthorizeTransactionUseCase{ur: acRepository, tr: tRepository, t: tr4}
	violations := useCase.Execute()

	if len(violations) == 2 {
		t.Errorf("lorem")
	}
}
