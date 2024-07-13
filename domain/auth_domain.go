package domain

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/messaging-api/types"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type AuthService interface {
	Login(ctx context.Context, payload types.LoginPayload) types.Response
}
