package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"

	"github.com/kieranajp/quiz/pkg/service"
)

type GameHandler struct {
	s *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{s: gameService}
}

func (gh *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	gameID, err := gh.s.CreateGame()
	if err != nil {
		http.Error(w, "Unable to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"game_id": gameID.String()})
}

func (gh *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "gameID")
	// Ensure it's a valid UUID
	gameID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	game, err := gh.s.GetGame(gameID)
	if err != nil {
		http.Error(w, "Unable to get game", http.StatusInternalServerError)
		return
	}
	err = gh.s.StartGame(game)
	if err != nil {
		http.Error(w, "Unable to start game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
