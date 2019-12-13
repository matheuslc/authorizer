package transactionentity

import (
	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// Transaction constants
const (
	AllowedAmountOfTransaction = 3
)

type TransactionValidation struct {
	User              ac.Account
	TransactionEvents <-chan es.Event
	CurrentEvent      es.Event
}

// AmountOfTransactions
func MoreThanAllowedViolation(tv TransactionValidation) bool {
	occurrences := []es.Event{}
	for event := range tv.TransactionEvents {
		occurrences = append(occurrences, event)
	}

	if len(occurrences) > AllowedAmountOfTransaction {
		return true
	}

	return false
}

// DuplicatedTransaction
func DuplicatedTransaction(tv TransactionValidation) (string, bool) {
	occurrences := []es.Event{}

	for event := range tv.TransactionEvents {
		sameMerchant := event.Payload.(Transaction).Merchant == tv.CurrentEvent.Payload.(Transaction).Merchant
		sameAmount := event.Payload.(Transaction).Amount == tv.CurrentEvent.Payload.(Transaction).Amount

		if sameMerchant && sameAmount {
			occurrences = append(occurrences, event)
		}
	}

	if len(occurrences) > 1 {
		return "doubled-transaction", true
	}

	return "", false
}

// AccountLimitViolation
func AccountLimitViolation(tv TransactionValidation) (string, bool) {
	balance := tv.User.AvailableLimit

	for event := range tv.TransactionEvents {
		balance += event.Payload.(Transaction).Amount
	}

	if tv.CurrentEvent.Payload.(Transaction).Amount < balance {
		return "insufficient-limit", true
	}

	return "", false
}

// AccountLimitViolation
func AccountNotInitilizedViolation(tv TransactionValidation) (string, bool) {
	nulla := ac.Account{}
	if tv.User == nulla {
		return "account-not-initialized", true
	}

	return "", false
}
