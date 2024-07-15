package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(rawTokenString string) string {
	return rawTokenString[7:]
}

func GetMapClaimsRaw(rawTokenString string) (jwt.MapClaims, error) {
	rawTokenString = ParseToken(rawTokenString)

	return GetMapClaims(rawTokenString)
}

func GetMapClaims(tokenString string) (jwt.MapClaims, error) {
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
