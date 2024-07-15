package repository

import (
	"context"

	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type messageRepitory struct {
	conn *config.PgxConfig
}

func NewMessageRepository(conn *config.PgxConfig) domain.MessageRepository {
	return &messageRepitory{
		conn: conn,
	}
}

func (m *messageRepitory) Create(ctx context.Context, payload types.CreateMessage) error {
	q := m.conn.TrOrDB(ctx)
	sql := `--sql
		INSERT INTO messages (id, chat_id, user_id, content) VALUES ($1, $2, $3, $4);
	`

	_, err := q.Exec(ctx, sql,
		payload.Id,
		payload.ChatId,
		payload.UserId,
		payload.Content,
	)

	return err
}
