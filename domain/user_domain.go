package domain

import (
	"context"

	"github.com/umardev500/messaging-api/types"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*types.User, error)
}
