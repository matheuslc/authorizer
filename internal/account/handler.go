package account

import (
	"github.com/matheuslc/authorizer/internal/eventstore"
	"github.com/matheuslc/authorizer/internal/eventstore/memorystore"
)

// AccountCommandHandler
func CommandHandler(
	payload map[string]interface{},
	acStore *memorystore.MemoryStore,
) eventstore.Event {
	acRepo := Repository{DB: acStore}

	useCase := CreateUseCase{
		AccountRepo: acRepo,
		AccountIntent: Account{
			ActiveCard:     payload["active-card"].(bool),
			AvailableLimit: int(payload["available-limit"].(float64)),
		},
	}

	return useCase.Execute()
}
