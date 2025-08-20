package service

import (
	"context"
	"errors"
	"io"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type StorageService interface {
	Store(ctx context.Context, file *models.File, content io.Reader) error
	Retrieve(ctx context.Context, file *models.File) (io.ReadCloser, error)
	ListUserFiles(ctx context.Context, userID string) ([]*models.File, error)
	DeleteFile(ctx context.Context, file *models.File) error
}

var (
	ErrFileNotFound = errors.New("file not found")
	ErrStorageFull  = errors.New("storage is full")
	ErrInvalidFile  = errors.New("invalid file")
)
