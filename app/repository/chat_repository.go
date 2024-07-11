package repository

import (
	"context"

	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type chatRepository struct {
	conn *config.PgxConfig
}

func NewChatRepository(conn *config.PgxConfig) domain.ChatRepository {
	return &chatRepository{
		conn: conn,
	}
}

func (cr *chatRepository) CreateChat(ctx context.Context, payload types.PushNewChatPayload) (err error) {
	q := cr.conn.TrOrDB(ctx)
	sql := `--sql
		INSERT INTO chats (id, chat_name) VALUES ($1, $2)
	`

	_, err = q.Exec(ctx, sql, payload.Room, payload.ChatName)

	return
}
