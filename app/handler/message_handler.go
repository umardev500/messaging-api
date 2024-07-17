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

func (m *messageHandler) GetMessage(c *fiber.Ctx) error {
	var room = c.Params("room")
	var msgType = c.Query("type")
	var date = c.Query("date")

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.Response{
			Code:    fiber.StatusBadRequest,
			Message: fiber.ErrBadRequest.Message,
			Error: &types.Error{
				Code: types.ValidationErr,
				Details: types.ErrDetail{
					Field:  "date",
					Filter: "=",
					Detail: "date must be in RFC3339 format",
				},
			},
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	resp := m.messageService.GetMessage(ctx, types.GetMessageParams{
		ChatId: room,
		Type:   types.GetMessageType(msgType),
		Date:   parsedDate,
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
