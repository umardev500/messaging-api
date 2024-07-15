package domain

import (
	"context"

	"github.com/umardev500/messaging-api/types"
)

type MessageRepository interface {
	Create(ctx context.Context, payload types.CreateMessage) error
}
