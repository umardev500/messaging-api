package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func UploadFile(file *multipart.FileHeader, uploadPath string) (location string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return
	}

	name := uploadPath + file.Filename
	err = os.WriteFile(name, content, 0644)
	if err != nil {
		return
	}

	return name, nil
}
