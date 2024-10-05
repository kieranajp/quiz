package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type QuestionAsked struct {
	ID         uuid.UUID `json:"event_id"`
	GameID     uuid.UUID `json:"game_id"`
	Round      int       `json:"round"`
	QuestionID uuid.UUID `json:"question_id"`
	Timestamp  time.Time `json:"created_at"`
}

func QuestionWasAsked(gameID uuid.UUID, round int, questionID uuid.UUID) *QuestionAsked {
	return &QuestionAsked{
		ID:         uuid.New(),
		GameID:     gameID,
		Round:      round,
		QuestionID: questionID,
		Timestamp:  time.Now(),
	}
}

func (e *QuestionAsked) EventID() uuid.UUID {
	return e.ID
}

func (e *QuestionAsked) EventType() string {
	return "question_was_asked"
}

func (e *QuestionAsked) AggregateID() uuid.UUID {
	return e.GameID
}

func (e *QuestionAsked) CreatedAt() time.Time {
	return e.Timestamp
}

func (e *QuestionAsked) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
