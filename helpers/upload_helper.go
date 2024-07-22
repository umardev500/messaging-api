package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/messaging-api/config"
)

type UploadResponse struct {
	Location string `json:"location"`
	Url      string `json:"url"`
}

func UploadFile(file *multipart.FileHeader, uploadPath string, updateFilename *string) (resp *UploadResponse, err error) {
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
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) && updateFilename == nil {
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return nil, err
		}

		log.Info().Msgf("Directory created: %s", uploadPath)
	}

	var fileName string = uuid.New().String() + filepath.Ext(file.Filename)
	var filePath string = uploadPath + fileName

	// if fileName is not nil that indicate this is update method
	if updateFilename != nil {
		filePath = uploadPath + *updateFilename
		fileName = *updateFilename
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return
	}

	return &UploadResponse{
		Location: filePath,
		Url:      config.GetConfig().Upload.StaticUrl + fileName,
	}, nil
}
