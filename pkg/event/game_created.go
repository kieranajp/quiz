package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type GameCreated struct {
	ID        uuid.UUID `json:"event_id"`
	GameID    uuid.UUID `json:"game_id"`
	Timestamp time.Time `json:"created_at"`
}

func GameWasCreated(gameID uuid.UUID) *GameCreated {
	return &GameCreated{
		ID:        uuid.New(),
		GameID:    gameID,
		Timestamp: time.Now(),
	}
}

func (e *GameCreated) EventID() uuid.UUID {
	return e.ID
}

func (e *GameCreated) EventType() string {
	return "game_created"
}

func (e *GameCreated) AggregateID() uuid.UUID {
	return e.GameID
}

func (e *GameCreated) CreatedAt() time.Time {
	return e.Timestamp
}

func (e *GameCreated) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
