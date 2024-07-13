package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func CheckAuth(c *fiber.Ctx) error {
	var ticket = uuid.New().String()

	var signingKey = []byte(os.Getenv("SECRET_KEY"))
	var tokenString = c.Get("Authorization")
	if tokenString == "" {
		return fiber.ErrUnauthorized
	}
	// Remove prefix
	tokenString = tokenString[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		log.Error().Msgf("error parsing token: %v | ticket: %s", err, ticket)
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		log.Error().Msgf("token is not map claims of jwt mapclaims: %v | ticket: %s", err, ticket)
		return fiber.ErrUnauthorized
	}

	c.Locals("user", claims["user"])
	return c.Next()
}
