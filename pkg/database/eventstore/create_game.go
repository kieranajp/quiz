package eventstore

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func CreateGame(db *sqlx.DB, gameID uuid.UUID) error {
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE %s (
			event_id UUID PRIMARY KEY,
			event_type TEXT,
			event_data JSONB,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`, gameTableName(gameID))

	_, err := db.Exec(createTableQuery)

	return err
}
