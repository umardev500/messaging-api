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
	GetChatList(c *fiber.Ctx) error
}

type ChatService interface {
	SaveMessage(ctx context.Context, data types.InputNewMessage)
	PushNewChat(ctx context.Context, payload types.CreateNewChatPayload) (types.Response, error)
	GetClaims(tokenString string) (types.Response, error)
	GetChatList(ctx context.Context, param types.GetChatListParam) types.Response
}

type ChatRepository interface {
	CreateChat(ctx context.Context, payload types.CreateNewChatPayload) error
	GetChatList(ctx context.Context, param types.GetChatListParam) ([]types.ChatList, error)
}
