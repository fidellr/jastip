package repository

import (
	"context"

	"github.com/fidellr/jastip/backend/plateu"
	"github.com/fidellr/jastip/backend/plateu/models"
)

type ImageRepository interface {
	StoreImage(ctx context.Context, m *models.Image) error
	FetchImages(ctx context.Context, filter *plateu.Filter) ([]*models.Image, string, error)
	GetImageByID(ctx context.Context, imageID string) (*models.Image, error)
	UpdateImageByID(ctx context.Context, imageID string, m *models.Image) error
	RemoveImageByID(ctx context.Context, imageID string) error
}
