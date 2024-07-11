package service

import (
	"context"

	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type chatService struct{}

func NewChatService() domain.ChatService {
	return &chatService{}
}

func (c *chatService) SaveMessage(ctx context.Context, data types.InputNewMessage) {}
