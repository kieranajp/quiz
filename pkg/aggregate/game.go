package aggregate

import (
	"github.com/google/uuid"
	"github.com/kieranajp/quiz/pkg/event"
)

type GameAggregate struct {
	GameID           uuid.UUID
	HasStarted       bool
	Players          int
	Settings         GameSettings
	CurrentRound     int
	CurrentRoundType string
}

type GameSettings struct {
	NumberOfRounds    int
	MaxPlayers        int
	QuestionsPerRound int
}

func NewGameAggregate(gameID uuid.UUID) GameAggregate {
	return GameAggregate{
		GameID:     gameID,
		HasStarted: false,
		Players:    0,
		Settings: GameSettings{
			NumberOfRounds:    3,
			MaxPlayers:        4,
			QuestionsPerRound: 5,
		},
	}
}

func (g *GameAggregate) ApplyGameCreated() {
	// No specific logic needed for creation
}

func (g *GameAggregate) ApplyPlayerJoined() {
	g.Players++
}

func (g *GameAggregate) ApplyGameStarted() {
	g.HasStarted = true
}

func (g *GameAggregate) ApplyRoundStarted(e event.RoundStarted) {
	g.CurrentRound = e.Round
	g.CurrentRoundType = e.RoundType
}

func (g *GameAggregate) IsFull() bool {
	return g.Players == g.Settings.MaxPlayers
}
