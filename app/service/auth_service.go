package service

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
)

type authService struct {
}

func NewAuthService() domain.AuthService {
	return &authService{}
}

func (s *authService) Login(ctx context.Context, payload types.LoginPayload) types.Response {
	var ticket = ctx.Value(types.ProcIdKey).(string)
	var resp types.Response = types.Response{
		Ticket: ticket,
	}

	// @Todo
	// Fetch user data from database
	//

	var key = []byte(os.Getenv("SECRET_KEY"))
	var user = types.UserClaim{
		Id:       "78901234-5678-9012-3456-789012345678",
		Username: "user2",
	}
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)

	if err != nil {
		resp.Code = fiber.StatusInternalServerError
		resp.Message = "Error signing token"
		return resp
	}

	resp.Message = "Login successful"
	resp.Data = map[string]interface{}{
		"token": signedToken,
	}

	return resp
}
