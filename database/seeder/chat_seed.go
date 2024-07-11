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

func (s *Seeder) ChatSeed(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	log.Info().Msg("📦 Seeding chats...")
	filePath := s.baseURL + "/chats_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Chat
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	copyCount, err := q.CopyFrom(
		ctx,
		pgx.Identifier{"chats"},
		[]string{"id", "chat_name", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			values := []any{
				rows[i].Id,
				rows[i].ChatName,
				rows[i].CreatedAt,
				rows[i].UpdatedAt,
			}

			return values, nil
		}),
	)

	time.Sleep(150 * time.Millisecond) // add delay
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error seeding chats 🚧")
		return
	}

	log.Info().Msgf("✅ Seeded %d chats", copyCount)

	return
}

func (s *Seeder) ChatDown(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	filePath := s.baseURL + "/chats_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Chat
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	var ids []string
	for _, row := range rows {
		ids = append(ids, row.Id)
	}

	sql := "DELETE FROM chats WHERE id = ANY($1)"
	_, err = q.Exec(ctx, sql, ids)

	time.Sleep(150 * time.Millisecond)
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error dropping chats 🚧")
		return
	}

	log.Info().Msgf("✅ Dropped %d chats", len(ids))

	return
}
