package eventstore

import (
	"fmt"
	"strings"

	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func gameTableName(gameID uuid.UUID) string {
	tableName := fmt.Sprintf("game_%s", gameID)
	return strings.ReplaceAll(tableName, "-", "_")
}

func RecordThat(db *sqlx.DB, e event.Event) error {
	eventData, err := e.ToJSON()
	if err != nil {
		return err
	}

	saveEventQuery := fmt.Sprintf(`
		INSERT INTO %s (event_id, event_type, event_data, created_at) VALUES ($1, $2, $3, $4)
	`, gameTableName(e.AggregateID()))

	_, err = db.Exec(saveEventQuery, e.EventID(), e.EventType(), eventData, e.CreatedAt())
	if err != nil {
		return err
	}

	return nil
}
