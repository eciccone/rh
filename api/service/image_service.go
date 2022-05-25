package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageService interface {
	SaveImage(file *multipart.FileHeader, path string, filename string) error
	DeleteImage(path string, filename string) error
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (s *imageService) SaveImage(file *multipart.FileHeader, path string, filename string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("SaveImage failed to open file: %w", err)
	}
	defer src.Close()

	dest, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return fmt.Errorf("SaveImage failed to create destination file: %w", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("SaveImage failed to copy source to destination: %w", err)
	}

	return nil
}

func (s *imageService) DeleteImage(path string, filename string) error {
	file := filepath.Join(path, filename)
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("DeleteImage failed to delete image file: %w", err)
	}

	return nil
}
