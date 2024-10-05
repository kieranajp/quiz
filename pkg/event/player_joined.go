package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type PlayerJoined struct {
	ID        uuid.UUID `json:"event_id"`
	GameID    uuid.UUID `json:"game_id"`
	PlayerID  uuid.UUID `json:"player_id"`
	Timestamp time.Time `json:"created_at"`
}

func PlayerHasJoined(gameID, playerID uuid.UUID) *PlayerJoined {
	return &PlayerJoined{
		ID:        uuid.New(),
		GameID:    gameID,
		PlayerID:  playerID,
		Timestamp: time.Now(),
	}
}

func (e *PlayerJoined) EventID() uuid.UUID {
	return e.ID
}

func (e *PlayerJoined) EventType() string {
	return "player_joined"
}

func (e *PlayerJoined) AggregateID() uuid.UUID {
	return e.GameID
}

func (e *PlayerJoined) CreatedAt() time.Time {
	return e.Timestamp
}

func (e *PlayerJoined) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
