package transactionentity

import (
	"time"
)

// TransactionValidated defines
const TransactionValidated = "transaction:validated"

// Transaction defines how a Transaction looks like
type Transaction struct {
	Merchant string
	Amount   int
	Time     time.Time
}
