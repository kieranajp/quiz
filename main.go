package main

import (
	"os"

	"github.com/kieranajp/quiz/cmd"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "db-dsn",
				Usage:   "Database DSN",
				Value:   "postgres://127.0.0.1:5432/quiz?sslmode=disable",
				EnvVars: []string{"DB_DSN"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "listen-addr",
						Usage:   "Listen Address",
						Value:   "0.0.0.0:8080",
						EnvVars: []string{"LISTEN_ADDRESS"},
					},
				},
				Action: cmd.Serve,
			},
		},
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Exit")
	}
}
