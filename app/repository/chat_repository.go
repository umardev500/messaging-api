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

func (cr *chatRepository) CreateChat(ctx context.Context, payload types.PushNewChatPayload) error {
	// q := cr.conn.TrOrDB(ctx)
	// sql := `--sql
	// 	INSERT INTO chat_participants(room_id, user_id) VALUES ($1, $2)
	// `

	return nil
}
