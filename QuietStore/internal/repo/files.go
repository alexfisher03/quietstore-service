package repo

import (
	"context"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type Files interface {
	Create(ctx context.Context, f *models.File) error
	ByID(ctx context.Context, id string) (*models.File, error)
	ListByOwner(ctx context.Context, ownerID string, limit, offset int) ([]*models.File, error)
	Delete(ctx context.Context, id string, ownerID string) error
}
