package transaction

import (
	"time"

	"github.com/google/uuid"
)

// TransactionValidated defines
const TransactionValidated = "transaction:validated"

// Transaction defines how a Transaction looks like
type Transaction struct {
	ID       uuid.UUID
	Merchant string
	Amount   int
	time     time.Time
}
