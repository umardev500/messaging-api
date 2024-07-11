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

func (s *Seeder) ChatParticipantSeed(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	log.Info().Msg("📦 Seeding chat participants...")
	filePath := s.baseURL + "/chat_participants.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Participant
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	copyCount, err := q.CopyFrom(
		ctx,
		pgx.Identifier{"chat_participants"},
		[]string{"chat_id", "user_id", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			values := []any{
				rows[i].ChatId,
				rows[i].UserId,
				rows[i].CreatedAt,
				rows[i].UpdatedAt,
			}

			return values, nil
		}),
	)

	time.Sleep(150 * time.Millisecond) // add delay
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error seeding chat participants 🚧")
		return
	}

	log.Info().Msgf("✅ Seeded %d chat participants", copyCount)

	return
}

func (s *Seeder) ChatParticipantDown(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	filePath := s.baseURL + "/chat_participants.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	var rows []types.Participant
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	var chatIds, userIds []string
	for _, row := range rows {
		chatIds = append(chatIds, row.ChatId)
		userIds = append(userIds, row.UserId)
	}

	sql := "DELETE FROM chat_participants WHERE chat_id = ANY($1) AND user_id = ANY($2)"
	_, err = q.Exec(ctx, sql, chatIds, userIds)

	time.Sleep(150 * time.Millisecond)
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error dropping chat participants 🚧")
		return
	}

	log.Info().Msgf("✅ Dropped %d chat participants", len(chatIds))

	return
}
