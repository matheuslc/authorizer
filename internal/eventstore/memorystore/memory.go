package memorystore

import (
	"sort"
	"sync"
	"time"

	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// MemoryStore defines how an event store need to be
type MemoryStore struct {
	sync.RWMutex
	Namespace string
	Data      map[string][]es.Event
}

// NewStorage is like a factory function to create a new storage
func NewStorage(namespace string) MemoryStore {
	data := make(map[string][]es.Event)
	return MemoryStore{Namespace: namespace, Data: data}
}

// Append a new event into the event store
func (db *MemoryStore) Append(event es.Event) es.Event {
	db.Lock()
	defer db.Unlock()

	currentHistory := db.Data[db.Namespace]
	newHistory := append(currentHistory, event)

	db.Data[db.Namespace] = newHistory
	// db.sort()
	return event
}

// Get the key value from store
func (db *MemoryStore) Get() []es.Event {
	db.Lock()
	defer db.Unlock()

	value := db.Data[db.Namespace]
	return value
}

// Iter iterates over the items in a concurrent map
func (db *MemoryStore) Iter() <-chan es.Event {
	c := make(chan es.Event)

	f := func() {
		db.Lock()
		defer db.Unlock()

		for _, event := range db.Data[db.Namespace] {
			c <- event
		}
		close(c)
	}
	go f()

	return c
}

// IterAfter iterates after a defined event
func (db *MemoryStore) IterAfter(after time.Time) <-chan es.Event {
	c := make(chan es.Event)
	data := db.Iter()

	f := func() {
		for event := range data {
			if event.Timestamp.After(after) || event.Timestamp.Equal(after) {
				c <- event
			}
		}

		close(c)
	}
	go f()

	return c
}

// EventsByName filters events by name and returns the collection
func (db *MemoryStore) EventsByName(name string) <-chan es.Event {
	c := make(chan es.Event)
	data := db.Iter()

	f := func() {
		for event := range data {
			if event.Name == name {
				c <- event
			}
		}

		close(c)
	}

	go f()
	return c
}

func (db *MemoryStore) sort() {
	sort.Slice(db.Data[db.Namespace], func(i, j int) bool {
		return db.Data[db.Namespace][i].Timestamp.Before(db.Data[db.Namespace][j].Timestamp)
	})
}
