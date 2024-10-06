package eventstore

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Aggregate defines the interface for all aggregates.
type Aggregate interface {
	ApplyEvent(event event.Event) error
	AggregateID() uuid.UUID
}

type BaseAggregate struct{}

// ApplyEvent dynamically calls an event handler method based on the event type.
func (b *BaseAggregate) ApplyEvent(target interface{}, e event.Event) error {
	// Create the method name by converting event type to a method name, e.g. "game_created" -> "ApplyGameCreated"
	c := cases.Title(language.English)
	methodName := "Apply" + strings.ReplaceAll(c.String(e.EventType()), "_", "")

	method := reflect.ValueOf(target).MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("unknown event type: %s", e.EventType())
	}

	// Invoke the method
	result := method.Call([]reflect.Value{reflect.ValueOf(e)})
	if len(result) == 1 && !result[0].IsNil() {
		return result[0].Interface().(error)
	}
	return nil
}
