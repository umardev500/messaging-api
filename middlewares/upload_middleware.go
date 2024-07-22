package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
	"github.com/umardev500/messaging-api/helpers"
	"github.com/umardev500/messaging-api/types"
)

func UpdateUploadMiddleware(uploadPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticket := uuid.New().String()
		var resp = types.Response{
			Ticket: ticket,
		}

		var dir = config.GetConfig().Upload.Path
		var filename = c.Params("filename")
		var fileLocation = dir + filename

		if info, err := os.Stat(fileLocation); os.IsNotExist(err) || info.IsDir() {
			resp.Message = fiber.ErrNotFound.Message
			return c.Status(fiber.ErrNotFound.Code).JSON(resp)
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			log.Err(err).Msgf("error getting file | err : %v", err)
			resp.Message = fiber.ErrBadRequest.Message

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		uploadResp, err := helpers.UploadFile(fileHeader, dir, &filename)
		if err != nil {
			log.Err(err).Msgf("error uploading file to server | err : %v", err)
			return err
		}

		resp.Message = "Replace upload"
		resp.Data = uploadResp
		return c.JSON(resp)
	}
}

func UploadMiddleware(uploadPath string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ticket := uuid.New().String()
		var resp = types.Response{
			Ticket: ticket,
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			log.Err(err).Msgf("error getting file | err : %v", err)
			resp.Message = fiber.ErrBadRequest.Message

			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}

		var dir = config.GetConfig().Upload.Path

		uploadResp, err := helpers.UploadFile(fileHeader, dir, nil)
		if err != nil {
			log.Err(err).Msgf("error uploading file to server | err : %v", err)
			return err
		}

		resp.Message = "Upload file"
		resp.Data = uploadResp
		return c.JSON(resp)
	}
}
