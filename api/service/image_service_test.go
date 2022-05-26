package service

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveImage(t *testing.T) {
	data := []struct {
		RemoveFile func(path string, filename string) error
		Assert     func(err error)
	}{
		{
			RemoveFile: func(path, filename string) error {
				return nil
			},
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			RemoveFile: func(path, filename string) error {
				return errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, d := range data {
		fp := FileProcessor{RemoveFile: d.RemoveFile}
		err := fp.DeleteImage("mockpath", "mock.file")
		d.Assert(err)
	}
}

func Test_SaveImage(t *testing.T) {
	data := []struct {
		OpenFile   func(file *multipart.FileHeader) (multipart.File, error)
		CreateFile func(path, filename string) (*os.File, error)
		CopyFile   func(w io.Writer, r io.Reader) (int64, error)
		Assert     func(err error)
	}{
		{
			OpenFile: func(file *multipart.FileHeader) (multipart.File, error) {
				return &os.File{}, nil
			},
			CreateFile: func(path, filename string) (*os.File, error) {
				return &os.File{}, nil
			},
			CopyFile: func(w io.Writer, r io.Reader) (int64, error) {
				return 0, nil
			},
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			OpenFile: func(file *multipart.FileHeader) (multipart.File, error) {
				return &os.File{}, errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			OpenFile: func(file *multipart.FileHeader) (multipart.File, error) {
				return &os.File{}, nil
			},
			CreateFile: func(path, filename string) (*os.File, error) {
				return &os.File{}, errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			OpenFile: func(file *multipart.FileHeader) (multipart.File, error) {
				return &os.File{}, nil
			},
			CreateFile: func(path, filename string) (*os.File, error) {
				return &os.File{}, nil
			},
			CopyFile: func(w io.Writer, r io.Reader) (int64, error) {
				return 0, errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, d := range data {
		fp := FileProcessor{OpenFile: d.OpenFile, CreateFile: d.CreateFile, CopyFile: d.CopyFile}
		err := fp.SaveImage(&multipart.FileHeader{}, "mockpath", "mock.file")
		d.Assert(err)
	}
}
