package eventstore

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event defines an event wrapper. Its payload contains the specific event
type Event struct {
	ID        uuid.UUID
	Timestamp time.Time
	Name      string
	Payload   interface{}
}

// EventStore interface defines what an EventStore needs
type EventStore interface {
	Append(event Event)
	Get() []Event
}

// InMemoryStorage defines how a storage is
type InMemoryStorage struct {
	sync.RWMutex
	Namespace string
	Data      map[string][]Event
}

// NewStorage is a factory to create a new storage
func NewStorage(namespace string) InMemoryStorage {
	data := make(map[string][]Event)
	return InMemoryStorage{Namespace: namespace, Data: data}
}

// Append a new event into the EventStore
func (db *InMemoryStorage) Append(event Event) {
	db.Lock()
	defer db.Unlock()

	currentHistory := db.Data[db.Namespace]
	newHistory := append(currentHistory, event)

	db.Data[db.Namespace] = newHistory
	db.sort()
}

// Get the key value from storage
func (db *InMemoryStorage) Get() ([]Event, bool) {
	db.Lock()
	defer db.Unlock()

	value, ok := db.Data[db.Namespace]
	return value, ok
}

func (db *InMemoryStorage) sort() {
	sort.Slice(db.Data[db.Namespace], func(i, j int) bool {
		return db.Data[db.Namespace][i].Timestamp.Before(db.Data[db.Namespace][j].Timestamp)
	})
}
