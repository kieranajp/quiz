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
}

// defaultEventStreamName is the default naming strategy, using a single table.
func defaultEventStreamName(_ uuid.UUID) string {
	return "event_stream"
}

// NewEventStore creates a new EventStore instance.
func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{
		subscribers:       make(map[string][]func(event interface{})),
		db:                db,
		eventStreamNaming: defaultEventStreamName,
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

func (es *EventStore) persistEvent(e event.Event) error {
	eventData, err := e.ToJSON()
	if err != nil {
		return err
	}

	saveEventQuery := fmt.Sprintf(`
		INSERT INTO %s (event_id, event_type, event_data, aggregate_id, created_at) VALUES ($1, $2, $3, $4, $5)
	`, es.eventStreamNaming(e.AggregateID()))

	_, err = es.db.Exec(saveEventQuery, e.EventID(), e.EventType(), eventData, e.AggregateID(), e.CreatedAt())
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
func (es *EventStore) LoadAggregate(aggregateID uuid.UUID, agg Aggregate) (Aggregate, error) {
	tableName := es.eventStreamNaming(aggregateID)

	query := fmt.Sprintf(
		"SELECT event_type, event_data FROM %s WHERE aggregate_id = $1 ORDER BY created_at ASC",
		tableName)

	var events []event.Event

	err := es.db.Select(&events, query, aggregateID)
	if err != nil {
		return nil, err
	}

	for _, e := range events {
		if err := agg.ApplyEvent(e); err != nil {
			return nil, err
		}
	}

	return agg, nil
}
