package service

import (
	"context"
	"io"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type StorageService interface {
	SaveFile(ctx context.Context, userID, originalName, contentType string, size int64, r io.Reader) (*models.File, error)
	OpenFile(ctx context.Context, userID, fileID string) (*models.File, io.ReadCloser, error)
	ListFiles(ctx context.Context, userID string, limit, offset int) ([]*models.File, error)
	DeleteFile(ctx context.Context, userID, fileID string) error
	SearchFiles(ctx context.Context, userID, q, contentType string, minSize, maxSize int64, limit, offset int) ([]*models.File, error)
	RenameFile(ctx context.Context, userID, fileID, newName string) error
}
