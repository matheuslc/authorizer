package account

import es "github.com/matheuslc/authorizer/internal/eventstore"

// Account violations string
const (
	AlreadyInitialized = "account-already-initialized"
	NotInitialized     = "account-not-initialized"
	CardNotActive      = "card-not-active"
	Empty              = ""
)

// Violations describes an Account violations structure to be validated
type Violations struct {
	AccountEvents []es.Event
	AccountIntent Account
}

// AlreadyInitializedViolation
func AlreadyInitializedViolation(v Violations) (string, bool) {
	for _, event := range v.AccountEvents {
		if event.Name == AccountCreated {
			return AlreadyInitialized, true
		}
	}

	return Empty, false
}

// NotInitilizedViolation checks if an account was previously initialized
func NotInitilizedViolation(v Violations) (string, bool) {
	_, ok := AlreadyInitializedViolation(v)

	if !ok {
		return NotInitialized, true
	}

	return Empty, false
}

// ActiveCardViolation checks if the account have an ActiveCard
func ActiveCardViolation(v Violations) (string, bool) {
	if v.AccountIntent.ActiveCard {
		return CardNotActive, false
	}

	return Empty, true
}
