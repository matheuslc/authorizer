package transaction

import (
	"time"
)

// TransactionValidated defines
const TransactionValidated = "transaction:validated"

// Transaction defines how a Transaction looks like
type Transaction struct {
	Merchant string `json:"merchant"`
	Amount   int    `json:"amount"`
	Time     time.Time
}
