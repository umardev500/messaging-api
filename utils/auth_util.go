package utils

import (
	"context"
	"fmt"
)

func GetUserIdFromLocals(ctx context.Context) (string, error) {
	userId := ctx.Value("user").(map[string]interface{})["id"]
	if userId == nil {
		return "", fmt.Errorf("no value")
	}

	return userId.(string), nil
}
