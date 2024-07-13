package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type authHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) domain.AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	var payload types.LoginPayload
	var ticket = uuid.New().String()
	var resp types.Response = types.Response{
		Ticket: ticket,
	}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, types.ProcIdKey, ticket)

	resp = a.authService.Login(ctx, payload)
	return c.JSON(resp)
}
