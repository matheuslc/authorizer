package transactionentity

import (
	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// Transaction constants
const (
	AllowedAmountOfTransaction = 3
)

// TransactionValidation
type TransactionValidation struct {
	User              ac.Account
	TransactionEvents <-chan es.Event
	CurrentEvent      es.Event
}

// MoreThanAllowedViolation checks if the account made more transactions than the allowed
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

// DuplicatedTransaction checks for duplicate transactions
func DuplicatedTransaction(tv TransactionValidation) (string, bool) {
	occurrences := []es.Event{}

	for event := range tv.TransactionEvents {
		sameMerchant := event.Payload.(Transaction).Merchant == tv.CurrentEvent.Payload.(Transaction).Merchant
		sameAmount := event.Payload.(Transaction).Amount == tv.CurrentEvent.Payload.(Transaction).Amount

		if sameMerchant && sameAmount {
			occurrences = append(occurrences, event)
		}
	}

	if len(occurrences) > 0 {
		return "doubled-transaction", true
	}

	return "", false
}

// AccountLimitViolation checks if the account have limit to continue with
// its current transaction
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

// AccountNotInitilizedViolation checks if an account was previously initialized.
// If not, the function will return false and "account-not-initialized"
func AccountNotInitilizedViolation(tv TransactionValidation) (string, bool) {
	nulla := ac.Account{}
	if tv.User == nulla {
		return "account-not-initialized", true
	}

	return "", false
}

// AccountActiveCardViolation checks if the account have an ActiveCard
func AccountActiveCardViolation(tv TransactionValidation) (string, bool) {
	if tv.User.ActiveCard {
		return "", false
	}

	return "card-not-active", true
}
