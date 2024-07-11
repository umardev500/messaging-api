package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

type chatService struct {
	chatRepository domain.ChatRepository
	conn           *config.PgxConfig
}

func NewChatService(chatRepository domain.ChatRepository, conn *config.PgxConfig) domain.ChatService {
	return &chatService{
		chatRepository: chatRepository,
		conn:           conn,
	}
}

func (c *chatService) SaveMessage(ctx context.Context, data types.InputNewMessage) {}

func (c *chatService) PushNewChat(ctx context.Context, payload types.PushNewChatPayload) {
	payload.Participants = append(payload.Participants, payload.UserId)
	participants, err := json.Marshal(payload.Participants)
	if err != nil {
		fmt.Println(err)
		return
	}

	var roomId = uuid.New().String() // replace with actual id returned after create new chat
	payload.Room = roomId

	err = c.conn.WithTransaction(ctx, func(ctx context.Context) (err error) {
		// First initialize the chat
		err = c.chatRepository.CreateChat(ctx, payload)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Initialize chat participants

		fmt.Println("ok")

		return
	})

	if err != nil {
		log.Error().Msgf("error creating chat: %v", err)
		return
	}

	if err := storage.Redis.Set(roomId, participants, 0); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("servie", payload.Participants)
}
