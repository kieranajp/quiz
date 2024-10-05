package service

import (
	"fmt"
	"github.com/kieranajp/quiz/pkg/database/entity"
	"github.com/kieranajp/quiz/pkg/database/query"

	"github.com/kieranajp/quiz/pkg/database/eventstore"
	"github.com/kieranajp/quiz/pkg/event"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GameServiceInterface interface {
	CreateGame() (uuid.UUID, error)
	AddPlayer(gameID uuid.UUID, playerID uuid.UUID) error
}

type GameService struct {
	db *sqlx.DB
}

func NewGameService(db *sqlx.DB) *GameService {
	return &GameService{db: db}
}

func (s *GameService) CreateGame() (uuid.UUID, error) {
	gameID := uuid.New()
	initialPlayerID := uuid.New()

	err := eventstore.CreateGame(s.db, gameID)
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

func (s *GameService) NextQuestion(gameID uuid.UUID) (entity.Question, error) {
	var question entity.Question

	// Load the current game aggregate
	game, err := eventstore.LoadGameAggregate(s.db, gameID)
	if err != nil {
		return question, err
	}

	if !game.HasStarted {
		return question, fmt.Errorf("game %s has not yet started", gameID)
	}

	question, err = query.GetQuestion(s.db, game.CurrentRoundType)
	if err != nil {
		return question, err
	}

	err = eventstore.RecordThat(s.db, event.QuestionWasAsked(gameID, game.CurrentRound, question.QuestionID))
	if err != nil {
		return question, err
	}

	return question, nil
}
