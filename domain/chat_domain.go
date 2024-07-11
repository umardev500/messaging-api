package domain

import "github.com/gofiber/fiber/v2"

type ChatHandler interface {
	WsChatList() fiber.Handler
	WsChat() fiber.Handler
	PushNewChat(c *fiber.Ctx) error
}
