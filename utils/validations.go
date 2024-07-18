package utils

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/messaging-api/types"
)

func ValidateDateResp(ctx context.Context, dateString string, c *fiber.Ctx) (handler error, err error) {
	_, err = time.Parse(time.RFC3339, dateString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.Response{
			Ticket:  ctx.Value(types.ProcIdKey).(string),
			Code:    fiber.StatusBadRequest,
			Message: fiber.ErrBadRequest.Message,
			Error: &types.Error{
				Code: types.ValidationErr,
				Details: types.ErrDetail{
					Field:  "date",
					Filter: "=",
					Detail: "date must be in RFC3339 format",
				},
			},
		}), err
	}

	return nil, err
}
