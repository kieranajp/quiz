package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type RoundStarted struct {
	ID        uuid.UUID `json:"event_id"`
	GameID    uuid.UUID `json:"game_id"`
	Round     int       `json:"round"`
	RoundType string    `json:"round_type"`
	Timestamp time.Time `json:"created_at"`
}

func RoundHasStarted(gameID uuid.UUID, round int, roundType string) *RoundStarted {
	return &RoundStarted{
		ID:        uuid.New(),
		GameID:    gameID,
		Round:     round,
		RoundType: roundType,
		Timestamp: time.Now(),
	}
}

func (e *RoundStarted) EventID() uuid.UUID {
	return e.ID
}

func (e *RoundStarted) EventType() string {
	return "round_started"
}

func (e *RoundStarted) AggregateID() uuid.UUID {
	return e.GameID
}

func (e *RoundStarted) CreatedAt() time.Time {
	return e.Timestamp
}

func (e *RoundStarted) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
