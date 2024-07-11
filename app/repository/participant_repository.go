package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type participantRepository struct {
	conn *config.PgxConfig
}

func NewParticipantRepository(conn *config.PgxConfig) domain.ParticipantRepository {
	return &participantRepository{
		conn: conn,
	}
}

func (pr *participantRepository) AddParticipant(ctx context.Context, payload types.InputParticipant) error {
	q := pr.conn.TrOrDB(ctx)

	copyCount, err := q.CopyFrom(
		ctx,
		pgx.Identifier{"chat_participants"},
		[]string{"chat_id", "user_id"},
		pgx.CopyFromSlice(len(payload.Participants), func(i int) ([]any, error) {
			values := []any{
				payload.ChatId,
				payload.Participants[i],
			}

			return values, nil
		}),
	)
	if err != nil {
		return err
	}
	if copyCount != int64(len(payload.Participants)) {
		return fmt.Errorf("expected %d rows inserted, got %d", len(payload.Participants), copyCount)
	}
	return nil
}
