package account

// AccountCreated defines the name of the fact when a new account was created
const AccountCreated = "account:created"

// Account defines how an Account looks like
type Account struct {
	ActiveCard     bool
	AvailableLimit int
}
