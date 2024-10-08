package eventstore

import (
	"encoding/json"
	"fmt"
	"github.com/kieranajp/quiz/pkg/event"
	"reflect"
)

// EventRegistry is used to manage event types and their corresponding struct types.
type EventRegistry struct {
	eventTypes map[string]reflect.Type
}

// NewEventRegistry creates a new EventRegistry.
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		eventTypes: make(map[string]reflect.Type),
	}
}

// RegisterEvent registers an event type with its corresponding struct type.
func (er *EventRegistry) RegisterEvent(eventStruct interface{}) {
	eventType := reflect.TypeOf(eventStruct).Elem().Name()
	er.eventTypes[eventType] = reflect.TypeOf(eventStruct).Elem()
}

// CreateEventInstance creates a new instance of an event based on the event type string.
func (er *EventRegistry) CreateEventInstance(eventType string, eventData []byte) (event.Event, error) {
	eventStructType, exists := er.eventTypes[eventType]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	// Create a new instance of the struct
	newEventPtr := reflect.New(eventStructType).Interface()

	// Unmarshal JSON data into the new event instance
	if err := json.Unmarshal(eventData, newEventPtr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return newEventPtr.(event.Event), nil
}
