package repo

import (
	"context"

	"github.com/alexfisher03/quietstore-service/QuietStore/internal/models"
)

type Users interface {
	Create(ctx context.Context, u *models.User) error
	ByID(ctx context.Context, id string) (*models.User, error)
	ByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}
