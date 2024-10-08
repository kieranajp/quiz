package cmd

import (
	"github.com/google/uuid"
	"github.com/kieranajp/quiz/pkg/database/eventstore"
	"github.com/kieranajp/quiz/pkg/event"
	"github.com/kieranajp/quiz/pkg/service"
	"net/http"
	"strings"

	"github.com/kieranajp/quiz/pkg/database"
	"github.com/kieranajp/quiz/pkg/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Serve(c *cli.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db, err := database.Connect(c.String("db-dsn"))
	if err != nil {
		return err
	}
	defer db.Close()

	eventRegistry := eventstore.NewEventRegistry()
	eventRegistry.RegisterEvent(&event.GameCreated{})
	eventRegistry.RegisterEvent(&event.GameStarted{})
	eventRegistry.RegisterEvent(&event.PlayerJoined{})
	eventRegistry.RegisterEvent(&event.RoundStarted{})
	eventRegistry.RegisterEvent(&event.QuestionAsked{})

	es := eventstore.NewEventStore(db, eventRegistry)
	es.SetEventStreamNamingStrategy(func(aggregateID uuid.UUID) string {
		return "game_" + strings.ReplaceAll(aggregateID.String(), "-", "_")
	})

	gs := service.NewGameService(es)
	gh := handler.NewGameHandler(gs)

	r.Get("/", handler.WelcomeHandler)
	r.Post("/game", gh.CreateGame)
	r.Post("/game/{gameID}/start", gh.StartGame)
	// r.Get("/game/{gameID}/question", gh.NextQuestion)

	log.Info().Str("listen-addr", c.String("listen-addr")).Msg("Starting server")
	return http.ListenAndServe(c.String("listen-addr"), r)
}
