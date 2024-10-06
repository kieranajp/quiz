package service

import (
	"fmt"
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
	gameID := uuid.New()
	initialPlayerID := uuid.New()

	err := s.es.CreateGame(s.db, gameID)
	if err != nil {
		return uuid.Nil, err
	}

	err = eventstore.RecordThat(s.db, event.GameWasCreated(gameID))
	if err != nil {
		return uuid.Nil, err
	}

	// Record the PlayerJoined event immediately after game creation
	err = s.AddPlayer(gameID, initialPlayerID)
	if err != nil {
		return uuid.Nil, err
	}

	return gameID, nil
}

func (s *GameService) AddPlayer(gameID uuid.UUID, playerID uuid.UUID) error {
	// Check if the game has started
	game, err := eventstore.LoadGameAggregate(s.db, gameID)
	if err != nil {
		return err
	}
	if game.HasStarted {
		return fmt.Errorf("cannot join game %s as it has already started", gameID)
	}
	if game.IsFull() {
		return fmt.Errorf("cannot join game %s as it is already full", gameID)
	}

	// Record the PlayerJoined event
	err = eventstore.RecordThat(s.db, event.PlayerHasJoined(gameID, playerID))
	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) StartGame(gameID uuid.UUID) error {
	// Load the current game aggregate
	game, err := eventstore.LoadGameAggregate(s.db, gameID)
	if err != nil {
		return err
	}

	if game.HasStarted {
		return fmt.Errorf("game %s has already started", gameID)
	}

	// Record the GameStarted event
	err = eventstore.RecordThat(s.db, event.GameHasStarted(gameID))
	if err != nil {
		return err
	}

	// The first round starts immediately
	err = eventstore.RecordThat(s.db, event.RoundHasStarted(gameID, 1, "multiple_choice"))
	if err != nil {
		return err
	}

	return nil
}
