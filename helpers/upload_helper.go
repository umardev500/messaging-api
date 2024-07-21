package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UploadResponse struct {
	Location string `json:"location"`
}

func UploadFile(file *multipart.FileHeader, uploadPath string, replaced *string) (location string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return
	}

	// Check for directory
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) && replaced == nil {
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return "", err
		}

		log.Info().Msgf("Directory created: %s", uploadPath)
	}

	var name string = uploadPath + uuid.New().String() + filepath.Ext(file.Filename)
	if replaced != nil {
		name = *replaced
	}

	err = os.WriteFile(name, content, 0644)
	if err != nil {
		return
	}

	return name, nil
}