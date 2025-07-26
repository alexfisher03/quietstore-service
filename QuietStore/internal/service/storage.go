package service

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type StorageService interface {
	Store(ctx context.Context, file *models.File, content io.Reader) error
	Retrieve(ctx context.Context, file *models.File) (io.ReadCloser, error)
	Delete(ctx context.Context, file *models.File) error
	GetPresignedURL(ctx context.Context, file *models.File, expiry time.Duration) (string, error)
}

var (
	ErrFileNotFound = errors.New("file not found")
	ErrStorageFull  = errors.New("storage is full")
	ErrInvalidFile  = errors.New("invalid file")
)
