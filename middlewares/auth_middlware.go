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
		log.Error().Msgf("error casting mapclaims: %v | ticket: %s", err, ticket)
		return fiber.ErrUnauthorized
	}

	c.Locals("user", claims["user"])
	return c.Next()
}

func GetMapClaims(tokenString string) (jwt.MapClaims, error) {
	tokenString = tokenString[7:] // remove the bearer prefix
	var signingKey = []byte(os.Getenv("SECRET_KEY"))

	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return nil, fmt.Errorf("error casting mapclaims: %v", err)
	}

	return claims, nil
}
