package transactionentity

import (
	ac "github.com/matheuslc/authorizer/internal/account"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// Transaction constants
const (
	AllowedAmountOfTransaction = 3
	MoreThanAllowed            = "high-frequency-small-interval"
	DoubledNotAllowed          = "doubled-transaction"
	InsuficientLimit           = "insufficient-limit"
	Empty                      = ""
)

// Violations describes a Transaction violations structure to be validated
type Violations struct {
	Account           ac.Account
	TransactionEvents []es.Event
	CurrentEvent      es.Event
}

// MoreThanAllowedViolation checks if the account made more transactions than the allowed
func MoreThanAllowedViolation(v Violations) (string, bool) {
	occurrences := []es.Event{}
	for _, event := range v.TransactionEvents {
		occurrences = append(occurrences, event)
	}

	if len(occurrences) >= AllowedAmountOfTransaction {
		return MoreThanAllowed, true
	}

	return Empty, false
}

// DuplicatedTransaction checks for duplicate transactions
func DuplicatedTransaction(v Violations) (string, bool) {
	occurrences := []es.Event{}

	for _, event := range v.TransactionEvents {
		sameMerchant := event.Payload.(Transaction).Merchant == v.CurrentEvent.Payload.(Transaction).Merchant
		sameAmount := event.Payload.(Transaction).Amount == v.CurrentEvent.Payload.(Transaction).Amount

		if sameMerchant && sameAmount {
			occurrences = append(occurrences, event)
		}
	}

	if len(occurrences) > 0 {
		return DoubledNotAllowed, true
	}

	return Empty, false
}

// AccountLimitViolation checks if the account have limit to continue with
// its current transaction
func AccountLimitViolation(v Violations) (string, bool) {
	balance := v.Account.AvailableLimit
	events := []es.Event{}

	for _, event := range v.TransactionEvents {
		events = append(events, event)
		balance += event.Payload.(Transaction).Amount
	}

	amount := v.CurrentEvent.Payload.(Transaction).Amount

	if amount > balance {
		return InsuficientLimit, true
	}

	return Empty, false
}