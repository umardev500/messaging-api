package service

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/domain"
	"github.com/umardev500/messaging-api/types"
	"github.com/umardev500/messaging-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository domain.UserRepository
}

func NewAuthService(userRepository domain.UserRepository) domain.AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Login(ctx context.Context, payload types.LoginPayload) types.Response {
	var ticket = ctx.Value(types.ProcIdKey).(string)
	var resp types.Response = types.Response{
		Ticket: ticket,
	}

	// Fetch user data from database
	user, err := s.userRepository.FindByUsername(ctx, payload.Username)

	if err != nil {
		log.Error().Msgf("error finding by username: %v | ticket: %s", err, ticket)

		if err == pgx.ErrNoRows {
			resp.Code = fiber.StatusNotFound
			resp.Message = fiber.ErrNotFound.Message
			return resp
		}

		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message
		return resp
	}

	// Matching password
	err = utils.ComparePassword(user.PasswordHash, payload.Password)

	if err != nil {

		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Error().Msgf("error password not match: | ticket: %s", ticket)
			resp.Code = fiber.StatusUnauthorized
			resp.Message = fiber.ErrUnauthorized.Message
			return resp
		}

		log.Error().Msgf("error comparing password: %v | ticket: %s", err, ticket)
		resp.Code = fiber.StatusInternalServerError
		resp.Message = fiber.ErrInternalServerError.Message

		return resp
	}

	// Pass claims to token
	var key = []byte(os.Getenv("SECRET_KEY"))
	var userClaim = types.UserClaim{
		Id:       user.ID,
		Username: user.Username,
	}
	claims := jwt.MapClaims{
		"user": userClaim,
		"exp":  time.Now().Add(time.Hour * 168).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(key)

	if err != nil {
		log.Error().Msgf("error signing token: %v | ticket: %s", err, ticket)
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
