package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
	"github.com/umardev500/messaging-api/utils"
)

type messageService struct {
	messageRepository domain.MessageRepository
}

func NewMessageService(messageRepository domain.MessageRepository) domain.MessageService {
	return &messageService{
		messageRepository: messageRepository,
	}
}

func (m *messageService) CreateMessage(ctx context.Context, payload types.CreateMessage) types.Response {
	var resp = types.Response{
		Ticket:  uuid.New().String(),
		Code:    fiber.StatusInternalServerError,
		Message: fiber.ErrInternalServerError.Message,
	}
	userId, err := utils.GetUserIdFromLocals(ctx)
	if err != nil {
		log.Error().Msgf("error when get user id from locals | err: %v | ticket: %s", err, resp.Ticket)
		return resp
	}

	id := uuid.New().String()
	payload.Id = id
	payload.UserId = userId

	if err := m.messageRepository.Create(ctx, payload); err != nil {
		log.Error().Msgf("error when creating message | err: %v | ticket: %s", err, resp.Ticket)
		return resp
	}

	resp.Code = fiber.StatusCreated
	resp.Message = "Success"
	return resp
}
