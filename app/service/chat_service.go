package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/storage"
	"github.com/umardev500/messaging-api/types"
)

type chatService struct{}

func NewChatService() domain.ChatService {
	return &chatService{}
}

func (c *chatService) SaveMessage(ctx context.Context, data types.InputNewMessage) {}

func (c *chatService) PushNewChat(ctx context.Context, payload types.PushNewChatPayload) {
	payload.Participants = append(payload.Participants, payload.UserId)
	participants, err := json.Marshal(payload.Participants)
	if err != nil {
		fmt.Println(err)
		return
	}

	// @Todo Create new chat to postgres storage
	//
	var roomId = "9" // replace with actual id returned after create new chat

	if err := storage.Redis.Set(roomId, participants, 0); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("servie", payload.Participants)
}
