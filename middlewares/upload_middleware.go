package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/helpers"
	"github.com/umardev500/messaging-api/types"
)

func UploadMiddleware(uploadPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticket := uuid.New().String()

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return err
		}

		var dir = config.GetConfig().Upload.Path

		location, err := helpers.UploadFile(fileHeader, dir)
		if err != nil {
			log.Err(err).Msgf("error uploading file to server | err : %v", err)
			return err
		}

		var resp = types.Response{
			Ticket: ticket,
			Data: helpers.UploadResponse{
				Location: location,
			},
		}
		return c.JSON(resp)
	}
}
