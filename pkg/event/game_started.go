package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type GameStarted struct {
	ID        uuid.UUID `json:"event_id"`
	GameID    uuid.UUID `json:"game_id"`
	Timestamp time.Time `json:"created_at"`
}

func GameHasStarted(gameID uuid.UUID) *GameStarted {
	return &GameStarted{
		ID:        uuid.New(),
		GameID:    gameID,
		Timestamp: time.Now(),
	}
}

func (e *GameStarted) EventID() uuid.UUID {
	return e.ID
}

func (e *GameStarted) EventType() string {
	return "game_started"
}

func (e *GameStarted) AggregateID() uuid.UUID {
	return e.GameID
}

func (e *GameStarted) CreatedAt() time.Time {
	return e.Timestamp
}

func (e *GameStarted) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
