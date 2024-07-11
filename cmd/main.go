package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/app"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msgf("error loading .env file: %v", err)
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	app := app.NewApp()
	if err := app.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("error running app")
	}
}
