package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
	"github.com/umardev500/messaging-api/utils"
)

type messageHandler struct {
	messageService domain.MessageService
}

func NewMessageHandler(messageService domain.MessageService) domain.MessageHandler {
	return &messageHandler{
		messageService: messageService,
	}
}

func (m *messageHandler) GetMessage(c *fiber.Ctx) error {
	var room = c.Params("room")
	var msgType = c.Query("type")
	var date = c.Query("date")
	var ticket = uuid.New().String()

	handler, err := utils.ValidateDateResp(date, c)
	if err != nil {
		log.Error().Msgf("date validation failed | err: %v | ticket: %s", err, ticket)
		return handler
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, types.ProcIdKey, ticket)

	resp := m.messageService.GetMessage(ctx, types.GetMessageParams{
		ChatId: room,
		Type:   types.GetMessageType(msgType),
		Date:   date,
	})

	return c.Status(resp.Code).JSON(resp)
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
