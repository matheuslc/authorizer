package memorystore

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	es "github.com/matheuslc/authorizer/internal/eventstore"
)

func TestAppend(t *testing.T) {
	namespace := "fake_name"
	MemoryStore := NewStorage(namespace)
	timeRange := []time.Time{}
	middleTime := time.Date(2019, 12, 13, 12, 0, 0, 0, time.UTC)

	timeRange = append(
		timeRange,
		time.Date(2019, 12, 11, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 12, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 14, 12, 0, 0, 0, time.UTC),
		middleTime,
	)

	for index, date := range timeRange {
		uuid, _ := uuid.NewUUID()
		fakeEvent := es.Event{ID: uuid, Timestamp: date, Name: uuid.String(), Payload: index}

		MemoryStore.Append(fakeEvent)
	}

	storageFinalState := MemoryStore.Get()
	expected := storageFinalState[2].Timestamp.Equal(middleTime)

	if expected == false {
		t.Errorf("Event out of order was inserted sorted")
	}
}

func TestIterAfter(t *testing.T) {
	namespace := "fake_name"
	memoryStore := NewStorage(namespace)

	timeRange := []time.Time{}
	middleTime := time.Date(2019, 12, 13, 12, 0, 0, 0, time.UTC)
	timeRange = append(
		timeRange,
		time.Date(2019, 12, 11, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 12, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 13, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 14, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 16, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 16, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 17, 12, 0, 0, 0, time.UTC),
		time.Date(2019, 12, 18, 12, 0, 0, 0, time.UTC),
	)

	for index, date := range timeRange {
		uuid, _ := uuid.NewUUID()
		fakeEvent := es.Event{ID: uuid, Timestamp: date, Name: strconv.Itoa(index), Payload: index}

		memoryStore.Append(fakeEvent)
	}

	items := memoryStore.IterAfter(middleTime)

	itemsFromChannel := []es.Event{}
	for event := range items {
		itemsFromChannel = append(itemsFromChannel, event)
	}

	if len(itemsFromChannel) != 7 {
		t.Errorf("InterAfter doest not only returned events after the Time.")
	}
}

func TestAppendConcurrent(t *testing.T) {
	namespace := "fake_name"
	MemoryStore := NewStorage(namespace)
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go appendAndMarkAsDone(&MemoryStore, &wg)
	}
	wg.Wait()

	items := MemoryStore.Get()

	if len(items) != 100 {
		t.Errorf("Expected 100 items to be appended into the storage.")
	}
}

func appendAndMarkAsDone(db *MemoryStore, wg *sync.WaitGroup) {
	uuid, _ := uuid.NewUUID()
	fakeEvent := es.Event{ID: uuid, Timestamp: time.Now(), Name: "lorem:ipsum", Payload: 10}
	db.Append(fakeEvent)
	wg.Done()
}
