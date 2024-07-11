package seeder

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/utils"
)

type Seeder struct {
	Conn    *config.PgxConfig
	Logger  *utils.Logger
	baseURL string
}

func NewSeeder(conn *config.PgxConfig) *Seeder {
	return &Seeder{
		Conn:    conn,
		Logger:  utils.NewLogger(),
		baseURL: "database/seeder/data",
	}
}

func (s *Seeder) Run(ctx context.Context) {
	err := s.Conn.WithTransaction(ctx, func(ctx context.Context) (err error) {
		// SEEDER DOWN
		err = s.MessageDown(ctx)
		if err != nil {
			return
		}

		err = s.ChatParticipantDown(ctx)
		if err != nil {
			return
		}

		err = s.UserDown(ctx)
		if err != nil {
			return
		}

		err = s.ChatDown(ctx)
		if err != nil {
			return
		}

		// SEEDER UP
		err = s.UserSeed(ctx)
		if err != nil {
			return
		}

		err = s.ChatSeed(ctx)
		if err != nil {
			return
		}

		err = s.ChatParticipantSeed(ctx)
		if err != nil {
			return
		}

		err = s.MessageSeed(ctx)
		if err != nil {
			return
		}

		return err
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error running seeder")
	}

}
