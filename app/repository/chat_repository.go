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

func (cr *chatRepository) GetChatList(ctx context.Context, param types.GetChatListParam) ([]types.ChatList, error) {
	q := cr.conn.TrOrDB(ctx)
	sql := `--sql
		SELECT c.id,
			CASE
				WHEN c.chat_name IS NOT NULL THEN c.chat_name
				ELSE (
					SELECT u.username
					FROM chat_participants cp2
						JOIN users u ON cp2.user_id = u.id
					WHERE cp2.chat_id = c.id
						AND u.id <> $1
					LIMIT 1
				)
			END AS chat_name,
			m.content as msg_higlight,
			m.created_at as last_msg_date
		FROM chats c
			JOIN chat_participants cp ON c.id = cp.chat_id
			LEFT JOIN LATERAL (
				SELECT m.content,
					m.created_at
				FROM messages m
				WHERE m.chat_id = c.id
				ORDER BY m.created_at DESC
				LIMIT 1
			) m ON TRUE
		WHERE cp.user_id = $2
			AND m.created_at > $3 -- the date is last date on last message on chat list
		ORDER BY m.created_at DESC
	`
	var chatList []types.ChatList = []types.ChatList{}

	rows, err := q.Query(ctx, sql, param.UserId, param.UserId, param.Date)
	if err != nil {
		return chatList, err
	}

	defer rows.Close()

	for rows.Next() {
		var each types.ChatList
		if err := rows.Scan(
			&each.Id,
			&each.ChatName,
			&each.Content,
			&each.LastMsgDate,
		); err != nil {
			return chatList, err
		}

		chatList = append(chatList, each)
	}

	return chatList, nil
}

func (cr *chatRepository) CreateChat(ctx context.Context, payload types.PushNewChatPayload) (err error) {
	q := cr.conn.TrOrDB(ctx)
	sql := `--sql
		INSERT INTO chats (id, chat_name) VALUES ($1, $2)
	`

	_, err = q.Exec(ctx, sql, payload.Room, payload.ChatName)

	return
}
