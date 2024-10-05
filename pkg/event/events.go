package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventID() uuid.UUID
	EventType() string
	AggregateID() uuid.UUID
	CreatedAt() time.Time
	ToJSON() ([]byte, error)
}
