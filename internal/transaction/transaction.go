package transaction

import (
	"time"

	"github.com/google/uuid"
)

// Transaction defines how a Transaction looks like
type Transaction struct {
	ID       uuid.UUID
	Merchant string
	Amount   int
	time     time.Timer
}
