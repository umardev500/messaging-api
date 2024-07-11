package seeder

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/types"
)

func (s *Seeder) MessageSeed(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	log.Info().Msg("📦 Seeding messages...")
	filePath := s.baseURL + "/message_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Message
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	copyCount, err := q.CopyFrom(
		ctx,
		pgx.Identifier{"messages"},
		[]string{"id", "chat_id", "user_id", "content", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			values := []any{
				rows[i].Id,
				rows[i].ChatId,
				rows[i].UserId,
				rows[i].Content,
				rows[i].CreatedAt,
				rows[i].UpdatedAt,
			}

			return values, nil
		}),
	)

	time.Sleep(150 * time.Millisecond) // add delay
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error seeding messages 🚧")
		return
	}

	log.Info().Msgf("✅ Seeded %d messages", copyCount)

	return
}

func (s *Seeder) MessageDown(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	filePath := s.baseURL + "/message_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Message
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	var ids []string
	for _, row := range rows {
		ids = append(ids, row.ChatId)
	}

	sql := "DELETE FROM messages WHERE chat_id = ANY($1)"
	_, err = q.Exec(ctx, sql, ids)

	time.Sleep(150 * time.Millisecond)
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error dropping messages 🚧")
		return
	}

	log.Info().Msgf("✅ Dropped %d messages", len(ids))

	return
}
