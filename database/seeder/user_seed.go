package seeder

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"password_hash"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func (s *Seeder) UserSeed(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	log.Info().Msg("📦 Seeding users...")

	filePath := s.baseURL + "/users_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rows []User
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	copyCount, err := q.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"id", "username", "password_hash", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			values := []any{
				rows[i].ID,
				rows[i].Username,
				rows[i].PasswordHash,
				rows[i].CreatedAt,
				rows[i].UpdatedAt,
			}
			return values, nil
		}),
	)

	time.Sleep(150 * time.Millisecond) // add delay
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error seeding users 🚧")
		return
	}

	log.Info().Msgf("✅ Seeded %d users", copyCount)

	return
}

func (s *Seeder) UserDown(ctx context.Context) (err error) {
	q := s.Conn.TrOrDB(ctx)
	log.Info().Msg("🔥 Dropping users...")

	filePath := s.baseURL + "/users_data.json"
	f, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var rows []User
	if err := json.Unmarshal(f, &rows); err != nil {
		return err
	}

	var ids []string
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	sql := "DELETE FROM users WHERE id = ANY($1)"
	_, err = q.Exec(ctx, sql, ids)

	time.Sleep(150 * time.Millisecond)
	s.Logger.UplineClearPrev()

	if err != nil {
		log.Err(err).Msg("error dropping users 🚧")
		return
	}

	log.Info().Msgf("✅ Dropped %d users", len(ids))

	return
}
