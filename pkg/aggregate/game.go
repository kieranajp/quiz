package aggregate

import (
	"github.com/kieranajp/quiz/pkg/database/eventstore"
	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
)

// GameAggregate represents the state of a game and implements the Aggregate interface.
type GameAggregate struct {
	eventstore.BaseAggregate
	ID               uuid.UUID
	RoundNumber      int
	Players          []uuid.UUID
	HasStarted       bool
	Settings         GameSettings
	CurrentRound     int
	CurrentRoundType string
}

type GameSettings struct {
	NumberOfRounds    int
	MaxPlayers        int
	QuestionsPerRound int
}

func NewGame() GameAggregate {
	gameID := uuid.New()
	return GameAggregate{
		ID:         gameID,
		HasStarted: false,
		Players:    []uuid.UUID{},
		Settings: GameSettings{
			NumberOfRounds:    3,
			MaxPlayers:        4,
			QuestionsPerRound: 5,
		},
	}
}

func (g *GameAggregate) AggregateID() uuid.UUID {
	return g.ID
}

func (g *GameAggregate) ApplyEvent(e event.Event) error {
	return g.BaseAggregate.Apply(g, e)
}

func (g *GameAggregate) ApplyGameCreated(e event.GameCreated) {
	// No specific logic needed for creation
}

func (g *GameAggregate) ApplyPlayerJoined(e event.PlayerJoined) {
	g.Players = append(g.Players, e.PlayerID)
}

func (g *GameAggregate) ApplyGameStarted(e event.GameStarted) {
	g.HasStarted = true
}

func (g *GameAggregate) ApplyRoundStarted(e event.RoundStarted) {
	g.CurrentRound = e.Round
	g.CurrentRoundType = e.RoundType
}

func (g *GameAggregate) IsFull() bool {
	return len(g.Players) == g.Settings.MaxPlayers
}
