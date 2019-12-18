package account

// AccountCreated defines the name of the fact when a new account was created
const (
	AccountCreated     = "account:created"
	AccountDenied      = "account:denied"
	AccountTransaction = "account:transaction"
)

// Account defines how an Account looks like
type Account struct {
	ActiveCard     bool `json:"active-card"`
	AvailableLimit int  `json:"available-limit"`
}
