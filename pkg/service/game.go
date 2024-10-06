package service

import (
	"fmt"
	"github.com/kieranajp/quiz/pkg/aggregate"
	"github.com/kieranajp/quiz/pkg/database/eventstore"
	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
)

type GameServiceInterface interface {
	CreateGame() (uuid.UUID, error)
	AddPlayer(gameID uuid.UUID, playerID uuid.UUID) error
}

type GameService struct {
	es *eventstore.EventStore
}

func NewGameService(es *eventstore.EventStore) *GameService {
	return &GameService{es: es}
}

func (s *GameService) CreateGame() (uuid.UUID, error) {
	initialPlayerID := uuid.New()

	game := aggregate.NewGame()
	err := s.es.CreateEventStream(game.AggregateID())
	if err != nil {
		return uuid.Nil, err
	}

	err = s.es.RecordThat(event.GameWasCreated(game.AggregateID()))
	if err != nil {
		return uuid.Nil, err
	}

	// Record the PlayerJoined event immediately after game creation
	err = s.AddPlayer(game, initialPlayerID)
	if err != nil {
		return uuid.Nil, err
	}

	return game.AggregateID(), nil
}

func (s *GameService) GetGame(gameID uuid.UUID) (aggregate.GameAggregate, error) {
	game := &aggregate.GameAggregate{}
	game.ID = gameID
	err := s.es.LoadAggregate(gameID, game)
	if err != nil {
		return aggregate.GameAggregate{}, err
	}

	return *game, nil
}

func (s *GameService) AddPlayer(game aggregate.GameAggregate, playerID uuid.UUID) error {
	if game.HasStarted {
		return fmt.Errorf("cannot join game %s as it has already started", game.AggregateID())
	}
	if game.IsFull() {
		return fmt.Errorf("cannot join game %s as it is already full", game.AggregateID())
	}

	// Record the PlayerJoined event
	err := s.es.RecordThat(event.PlayerHasJoined(game.AggregateID(), playerID))
	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) StartGame(game aggregate.GameAggregate) error {
	if game.HasStarted {
		return fmt.Errorf("game %s has already started", game.AggregateID())
	}

	// Record the GameStarted event
	err := s.es.RecordThat(event.GameHasStarted(game.AggregateID()))
	if err != nil {
		return err
	}

	// The first round starts immediately
	err = s.es.RecordThat(event.RoundHasStarted(game.AggregateID(), 1, "multiple_choice"))
	if err != nil {
		return err
	}

	return nil
}
