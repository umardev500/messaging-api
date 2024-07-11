package service

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

type chatService struct {
	chatRepository        domain.ChatRepository
	participantRepository domain.ParticipantRepository
	conn                  *config.PgxConfig
}

func NewChatService(chatRepository domain.ChatRepository, participantRepository domain.ParticipantRepository, conn *config.PgxConfig) domain.ChatService {
	return &chatService{
		chatRepository:        chatRepository,
		participantRepository: participantRepository,
		conn:                  conn,
	}
}

func (c *chatService) SaveMessage(ctx context.Context, data types.InputNewMessage) {}

func (c *chatService) PushNewChat(ctx context.Context, payload types.PushNewChatPayload) (resp types.Response, err error) {
	ticket := uuid.New().String()
	resp.Ticket = ticket

	payload.Participants = append(payload.Participants, payload.UserId)
	participants, err := json.Marshal(payload.Participants)
	if err != nil {
		log.Error().Msgf("error marshaling participants: %v | ticket: %s", err, ticket)
		return
	}

	var roomId = uuid.New().String() // replace with actual id returned after create new chat
	payload.Room = roomId

	err = c.conn.WithTransaction(ctx, func(ctx context.Context) (err error) {
		// First initialize the chat
		err = c.chatRepository.CreateChat(ctx, payload)
		if err != nil {
			return
		}

		// Initialize chat participants
		err = c.participantRepository.AddParticipant(ctx, types.InputParticipant{
			ChatId:       roomId,
			Participants: payload.Participants,
		})

		return
	})

	if err != nil {
		log.Error().Msgf("failed at transaction: %v | ticket: %s", err, ticket)
		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message
		return
	}

	if err = storage.Redis.Set(roomId, participants, 0); err != nil {
		log.Error().Msgf("error setting redis: %v | ticket: %s", err, ticket)
		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message
		return
	}

	resp.Code = fiber.StatusCreated
	resp.Message = "Success create new chat"

	return
}
