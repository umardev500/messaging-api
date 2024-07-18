package domain

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/messaging-api/types"
)

type ChatHandler interface {
	WsChatList() fiber.Handler
	WsChat() fiber.Handler
	PushNewChat(c *fiber.Ctx) error
}

type ChatService interface {
	SaveMessage(ctx context.Context, data types.InputNewMessage)
	PushNewChat(ctx context.Context, payload types.PushNewChatPayload) (types.Response, error)
	GetClaims(tokenString string) (types.Response, error)
}

type ChatRepository interface {
	CreateChat(ctx context.Context, payload types.PushNewChatPayload) error
	GetChatList(ctx context.Context, param types.GetChatListParam) ([]types.ChatList, error)
}
