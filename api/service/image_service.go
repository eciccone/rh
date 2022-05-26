package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileProcessor struct {
	OpenFile   func(file *multipart.FileHeader) (multipart.File, error)
	CreateFile func(path string, filename string) (*os.File, error)
	CopyFile   func(w io.Writer, r io.Reader) (int64, error)
	RemoveFile func(path string, filename string) error
}

type ImageService interface {
	SaveImage(file *multipart.FileHeader, path string, filename string) error
	DeleteImage(path string, filename string) error
}

func NewImageService() ImageService {
	return &FileProcessor{
		OpenFile: func(file *multipart.FileHeader) (multipart.File, error) {
			return file.Open()
		},

		CreateFile: func(path, filename string) (*os.File, error) {
			return os.Create(filepath.Join(path, filename))
		},

		CopyFile: func(w io.Writer, r io.Reader) (int64, error) {
			return io.Copy(w, r)
		},

		RemoveFile: func(path string, filename string) error {
			return os.Remove(filepath.Join(path, filename))
		},
	}
}

func (f *FileProcessor) SaveImage(file *multipart.FileHeader, path string, filename string) error {
	src, err := f.OpenFile(file)
	if err != nil {
		return fmt.Errorf("SaveImage failed to open file: %w", err)
	}
	defer src.Close()

	dest, err := f.CreateFile(path, filename)
	if err != nil {
		return fmt.Errorf("SaveImage failed to create destination file: %w", err)
	}
	defer dest.Close()

	_, err = f.CopyFile(dest, src)
	if err != nil {
		return fmt.Errorf("SaveImage failed to copy source to destination: %w", err)
	}

	return nil
}

func (f *FileProcessor) DeleteImage(path string, filename string) error {
	err := f.RemoveFile(path, filename)
	if err != nil {
		return fmt.Errorf("DeleteImage failed to delete image file: %w", err)
	}

	return nil
}
