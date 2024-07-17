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

func (m *messageRepitory) GetMessage(ctx context.Context, params types.GetMessageParams) ([]types.Message, error) {
	q := m.conn.TrOrDB(ctx)
	var sql string

	if params.Type == types.MessageDown {
		sql = `--sql
		SELECT id, chat_id, user_id, content, created_at, updated_at
		FROM messages
		WHERE chat_id = $1 AND created_at > $2
		ORDER BY created_at
		LIMIT 10
	`
	} else {
		sql = `--sql
		SELECT id, chat_id, user_id, content, created_at, updated_at
		FROM messages
		WHERE chat_id = $1 AND created_at < $2
		ORDER BY created_at
		LIMIT 10
	`
	}

	rows, err := q.Query(ctx, sql, params.ChatId, params.Date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rowsData []types.Message
	for rows.Next() {
		var row types.Message
		if err := rows.Scan(&row.Id, &row.ChatId, &row.UserId, &row.Content, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		rowsData = append(rowsData, row)
	}

	return rowsData, nil
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
