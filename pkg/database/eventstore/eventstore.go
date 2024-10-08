package eventstore

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// EventStreamNamingStrategy defines a function type for naming event streams.
type EventStreamNamingStrategy func(aggregateID uuid.UUID) string

// EventStore manages event recording and subscriptions.
type EventStore struct {
	subscribers       map[string][]func(event interface{})
	mux               sync.Mutex
	db                *sqlx.DB
	eventStreamNaming EventStreamNamingStrategy
	eventRegistry     *EventRegistry
}

// defaultEventStreamName is the default naming strategy, using a single table.
func defaultEventStreamName(_ uuid.UUID) string {
	return "event_stream"
}

// NewEventStore creates a new EventStore instance.
func NewEventStore(db *sqlx.DB, registry *EventRegistry) *EventStore {
	return &EventStore{
		subscribers:       make(map[string][]func(event interface{})),
		db:                db,
		eventStreamNaming: defaultEventStreamName,
		eventRegistry:     registry,
	}
}

// SetEventStreamNamingStrategy allows setting a custom strategy for event stream naming.
func (es *EventStore) SetEventStreamNamingStrategy(strategy EventStreamNamingStrategy) {
	es.mux.Lock()
	defer es.mux.Unlock()
	es.eventStreamNaming = strategy
}

func (es *EventStore) CreateEventStream(aggregateID uuid.UUID) error {
	tableName := es.eventStreamNaming(aggregateID)

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			event_id UUID PRIMARY KEY,
			event_type TEXT,
			event_data JSONB,
			aggregate_id UUID,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_aggregate_id ON %s (aggregate_id);
	`, tableName, tableName)

	_, err := es.db.Exec(createTableQuery)

	return err
}

// RecordThat records an event, persists it, and notifies the subscribers.
func (es *EventStore) RecordThat(event event.Event) error {
	es.mux.Lock()
	defer es.mux.Unlock()

	err := es.persistEvent(event)
	if err != nil {
		return fmt.Errorf("failed to persist event: %w", err)
	}

	es.notifySubscribers(event)

	return nil
}

func (es *EventStore) eventType(e event.Event) string {
	return reflect.TypeOf(e).Elem().Name()
}

func (es *EventStore) persistEvent(e event.Event) error {
	eventData, err := e.ToJSON()
	if err != nil {
		return err
	}

	saveEventQuery := fmt.Sprintf(`
		INSERT INTO %s (event_id, event_type, event_data, aggregate_id, created_at) VALUES ($1, $2, $3, $4, $5)
	`, es.eventStreamNaming(e.AggregateID()))

	_, err = es.db.Exec(saveEventQuery, e.EventID(), es.eventType(e), eventData, e.AggregateID(), e.CreatedAt())
	if err != nil {
		return err
	}

	return nil
}

// notifySubscribers notifies all subscribers of a given event.
func (es *EventStore) notifySubscribers(event interface{}) {
	eventType := reflect.TypeOf(event).Name()
	if handlers, ok := es.subscribers[eventType]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// Subscribe allows registering a handler for a specific event type.
func (es *EventStore) Subscribe(eventType string, handler func(event interface{})) {
	es.mux.Lock()
	defer es.mux.Unlock()

	if _, ok := es.subscribers[eventType]; !ok {
		es.subscribers[eventType] = []func(event interface{}){}
	}
	es.subscribers[eventType] = append(es.subscribers[eventType], handler)
}

// LoadAggregate loads the aggregate state from the stored events.
func (es *EventStore) LoadAggregate(aggregateID uuid.UUID, agg Aggregate) error {
	tableName := es.eventStreamNaming(aggregateID)

	query := fmt.Sprintf(
		"SELECT event_type, event_data FROM %s WHERE aggregate_id = $1 ORDER BY created_at ASC",
		tableName)

	// Define a slice of structs to hold the raw event data
	var rawEvents []struct {
		EventType string `db:"event_type"`
		EventData []byte `db:"event_data"`
	}

	// Query the database to load the raw event data
	err := es.db.Select(&rawEvents, query, aggregateID)
	if err != nil {
		return err
	}

	// Iterate through the raw events and use the EventRegistry to create event instances
	for _, rawEvent := range rawEvents {
		e, err := es.eventRegistry.CreateEventInstance(rawEvent.EventType, rawEvent.EventData)
		if err != nil {
			return err
		}

		// Apply the event to the aggregate
		if err := agg.ApplyEvent(e); err != nil {
			return err
		}
	}

	return nil
}
