package repository

import (
	"context"

	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type userRepository struct {
	conn *config.PgxConfig
}

func NewuserRepository(conn *config.PgxConfig) domain.UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (u *userRepository) FindByUsername(ctx context.Context, username string) (user *types.User, err error) {
	q := u.conn.TrOrDB(ctx)
	sql := `--sql
		SELECT u.id, u.username FROM users u WHERE username=$1;
	`
	var data types.User

	if err := q.QueryRow(ctx, sql, username).Scan(&data.ID, &data.Username); err != nil {
		return nil, err
	}
	user = &data

	return
}
