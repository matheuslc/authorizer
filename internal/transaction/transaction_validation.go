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
	TransactionEvents []es.Event
	CurrentEvent      es.Event
}

// MoreThanAllowedViolation checks if the account made more transactions than the allowed
func MoreThanAllowedViolation(tv TransactionValidation) (string, bool) {
	occurrences := []es.Event{}
	for _, event := range tv.TransactionEvents {
		occurrences = append(occurrences, event)
	}

	if len(occurrences) >= AllowedAmountOfTransaction {
		return "high-frequency-small-interval", true
	}

	return "", false
}

// DuplicatedTransaction checks for duplicate transactions
func DuplicatedTransaction(tv TransactionValidation) (string, bool) {
	occurrences := []es.Event{}

	for _, event := range tv.TransactionEvents {
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
	events := []es.Event{}

	for _, event := range tv.TransactionEvents {
		events = append(events, event)
		balance += event.Payload.(Transaction).Amount
	}

	amount := tv.CurrentEvent.Payload.(Transaction).Amount

	if amount > balance {
		return "insufficient-limit", true
	}

	return "", false
}
