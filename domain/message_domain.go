package domain

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/messaging-api/types"
)

type MessageHandler interface {
	Create(c *fiber.Ctx) error
}

type MessageService interface {
	CreateMessage(ctx context.Context, payload types.CreateMessage) types.Response
	GetMessage(ctx context.Context, params types.GetMessageParams) types.Response
}

type MessageRepository interface {
	Create(ctx context.Context, payload types.CreateMessage) error
	GetMessage(ctx context.Context, params types.GetMessageParams) ([]types.Message, error)
}
