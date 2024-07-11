package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/realtime-chat/config"
	"github.com/umardev500/realtime-chat/database/seeder"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msgf("error loading .env file: %v", err)
	}
}

func main() {
	pgxConn := config.NewPgx()
	s := seeder.NewSeeder(pgxConn)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	s.Run(ctx)
}
