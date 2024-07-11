package domain

import (
	"context"

	"github.com/umardev500/messaging-api/types"
)

type ParticipantRepository interface {
	AddParticipant(ctx context.Context, payload types.InputParticipant) error
}
