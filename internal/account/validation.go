package account

// Account violations string
const (
	NotInitialized = "account-not-initialized"
	CardNotActive  = "card-not-active"
	Empty          = ""
)

// Violations describes an Account violations structure to be validated
type Violations struct {
	Account Account
}

// NotInitilizedViolation checks if an account was previously initialized
func NotInitilizedViolation(v Violations) (string, bool) {
	nulla := Account{}
	if v.Account == nulla {
		return NotInitialized, true
	}

	return Empty, false
}

// ActiveCardViolation checks if the account have an ActiveCard
func ActiveCardViolation(v Violations) (string, bool) {
	if v.Account.ActiveCard {
		return CardNotActive, false
	}

	return Empty, true
}
