package account

import "github.com/google/uuid"

// Account defines how an Account looks like
type Account struct {
	ID             uuid.UUID
	ActiveCard     bool
	AvailableLimit int
}
