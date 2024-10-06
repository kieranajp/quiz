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

// Apply dynamically calls an event handler method based on the event type.
func (b *BaseAggregate) Apply(target interface{}, e event.Event) error {
	// Create the method name by converting event type to a method name, e.g. "game_created" -> "ApplyGameCreated"
	c := cases.Title(language.English)
	snakeToCamel := func(input string) string {
		words := strings.Split(input, "_")
		for i := range words {
			words[i] = c.String(words[i])
		}
		return strings.Join(words, "")
	}

	methodName := "Apply" + snakeToCamel(e.EventType())

	// Use reflection to find the method on the target
	targetValue := reflect.ValueOf(target)
	method := targetValue.MethodByName(methodName)
	if !method.IsValid() {
		return fmt.Errorf("unknown event type: %s", e.EventType())
	}

	// Prepare the argument value
	eventValue := reflect.ValueOf(e)
	if eventValue.Kind() == reflect.Ptr {
		eventValue = reflect.Indirect(eventValue)
	}

	// Invoke the method
	result := method.Call([]reflect.Value{eventValue})
	if len(result) == 1 && !result[0].IsNil() {
		return result[0].Interface().(error)
	}
	return nil
}
