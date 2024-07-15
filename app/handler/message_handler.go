package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type messageHandler struct {
	messageService domain.MessageService
}

func NewMessageHandler(messageService domain.MessageService) domain.MessageHandler {
	return &messageHandler{
		messageService: messageService,
	}
}

func (m *messageHandler) Create(c *fiber.Ctx) error {
	var payload types.CreateMessage
	var ticket = uuid.New().String()
	var resp = types.Response{
		Ticket: ticket,
	}
	var room = c.Params("room")
	payload.ChatId = room

	if err := c.BodyParser(&payload); err != nil {
		resp.Message = fiber.ErrUnprocessableEntity.Message
		return c.Status(fiber.StatusUnprocessableEntity).JSON(resp)
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, types.ProcIdKey, ticket)

	resp = m.messageService.CreateMessage(ctx, payload)

	return c.Status(resp.Code).JSON(resp)
}
