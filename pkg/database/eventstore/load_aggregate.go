package eventstore

import (
	"encoding/json"
	"fmt"

	"github.com/kieranajp/quiz/pkg/aggregate"
	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func LoadGameAggregate(db *sqlx.DB, gameID uuid.UUID) (aggregate.GameAggregate, error) {
	tableName := gameTableName(gameID)

	query := fmt.Sprintf(
		"SELECT event_type, event_data FROM %s ORDER BY created_at ASC",
		tableName)

	events := []struct {
		EventType string `db:"event_type"`
		EventData []byte `db:"event_data"`
	}{}

	err := db.Select(&events, query)
	if err != nil {
		return aggregate.GameAggregate{}, err
	}

	game := aggregate.NewGameAggregate(gameID)
	for _, e := range events {
		switch e.EventType {
		case "game_created":
			game.ApplyGameCreated()
		case "player_joined":
			game.ApplyPlayerJoined()
		case "game_started":
			game.ApplyGameStarted()
		case "round_started":
			var roundStartedEvent event.RoundStarted
			err := json.Unmarshal(e.EventData, &roundStartedEvent)
			if err != nil {
				return aggregate.GameAggregate{}, err
			}
			game.ApplyRoundStarted(roundStartedEvent)
		}
	}

	return game, nil
}
