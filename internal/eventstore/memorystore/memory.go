package memorystore

import (
	"sort"
	"sync"
	"time"

	es "github.com/matheuslc/authorizer/internal/eventstore"
)

// MemoryStore defines how a storage is
type MemoryStore struct {
	sync.RWMutex
	Namespace string
	Data      map[string][]es.Event
}

// NewStorage is a factory to create a new storage
func NewStorage(namespace string) MemoryStore {
	data := make(map[string][]es.Event)
	return MemoryStore{Namespace: namespace, Data: data}
}

// Append a new event into the EventStore
func (db *MemoryStore) Append(event es.Event) es.Event {
	db.Lock()
	defer db.Unlock()

	currentHistory := db.Data[db.Namespace]
	newHistory := append(currentHistory, event)

	db.Data[db.Namespace] = db.sort(newHistory)

	return event
}

// Get the key value from storage
func (db *MemoryStore) Get() []es.Event {
	db.Lock()
	defer db.Unlock()

	value := db.Data[db.Namespace]
	return value
}

func (db *MemoryStore) sort(data []es.Event) []es.Event {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Before(db.Data[db.Namespace][j].Timestamp)
	})

	return data
}

// Iter over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword locking the resource
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

// EventsAfter collects an amount of data
func (db *MemoryStore) EventsAfter(after time.Time) []es.Event {
	db.Lock()
	defer db.Unlock()

	data := db.Data[db.Namespace]
	filteredEvents := []es.Event{}

	for _, event := range data {
		if event.Timestamp.After(after) || event.Timestamp.Equal(after) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}
