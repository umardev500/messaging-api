package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
	"github.com/umardev500/messaging-api/utils"
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

func (c *chatService) GetChatList(ctx context.Context, param types.GetChatListParam) types.Response {
	var resp types.Response = types.Response{
		Ticket:  ctx.Value(types.ProcIdKey).(string),
		Code:    fiber.ErrInternalServerError.Code,
		Message: fiber.ErrInternalServerError.Message,
	}

	cl, err := c.chatRepository.GetChatList(ctx, param)
	if err != nil {
		log.Error().Msgf("error when get chatlist | err: %v | ticket: %s", err, resp.Ticket)
		return resp
	}

	resp.Data = cl
	resp.Code = fiber.StatusOK
	resp.Message = "Get chat list"
	return resp
}

func (c *chatService) GetClaims(tokenString string) (types.Response, error) {
	var ticket = uuid.New().String()

	var resp = types.Response{
		Ticket: ticket,
	}

	// Validating
	if tokenString == "" {
		log.Error().Msgf("token string is empy | ticket: %s", ticket)
		resp.Message = fiber.ErrBadRequest.Message
		resp.Error = &types.Error{
			Code: types.ValidationErr,
			Details: types.ErrDetail{
				Field:  "token",
				Filter: "required",
				Detail: "token parameter is required",
			},
		}

		return resp, fmt.Errorf("token string is empty")
	}
	claims, err := utils.GetMapClaims(tokenString)
	if err != nil {
		log.Error().Msgf("failed to get map claims err: %v | ticket: %s", err, ticket)
		resp.Message = fiber.ErrBadRequest.Message
		resp.Error = &types.Error{
			Code: types.ValidationErr,
			Details: types.ErrDetail{
				Field:  "token",
				Filter: "failed",
				Detail: "token is failed on check",
			},
		}

		return resp, fmt.Errorf("token is failed on check")
	}

	resp.Data = claims

	return resp, nil
}

func (c *chatService) SaveMessage(ctx context.Context, data types.InputNewMessage) {}

func (c *chatService) PushNewChat(ctx context.Context, payload types.CreateNewChatPayload) (resp types.Response, err error) {
	timestamp := time.Now().UTC().Unix()
	ticket := uuid.New().String()
	resp.Ticket = ticket

	payload.Participants = append(payload.Participants, payload.UserId)

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

	// Set redis
	type RedisValue struct {
		Blocked bool `json:"blocked"`
	}

	value := map[string]interface{}{}
	for _, participant := range payload.Participants {
		value[participant] = RedisValue{Blocked: false}
	}

	valueBytes, err := json.Marshal(value)

	if err != nil {
		log.Error().Msgf("failed to marshal value: %v | ticket: %s", err, ticket)
		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message
		return
	}

	if err = storage.Redis.JSONSet(ctx, roomId, "$", valueBytes).Err(); err != nil {
		log.Error().Msgf("error setting redis: %v | ticket: %s", err, ticket)
		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message
		return
	}

	resp.Code = fiber.StatusCreated
	resp.Message = "Success create new chat"
	resp.Data = types.BroadcastChatList{
		Room:      roomId,
		Message:   *payload.Message.Text,
		Timestamp: timestamp,
	}

	return
}
