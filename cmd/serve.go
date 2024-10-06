package cmd

import (
	"github.com/kieranajp/quiz/pkg/database/eventstore"
	"github.com/kieranajp/quiz/pkg/service"
	"net/http"

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

	es := eventstore.NewEventStore(db)

	gs := service.NewGameService(es)
	gh := handler.NewGameHandler(gs)

	r.Get("/", handler.WelcomeHandler)
	r.Post("/game", gh.CreateGame)
	r.Post("/game/{gameID}/start", gh.StartGame)
	r.Get("/game/{gameID}/question", gh.NextQuestion)

	log.Info().Str("listen-addr", c.String("listen-addr")).Msg("Starting server")
	return http.ListenAndServe(c.String("listen-addr"), r)
}
