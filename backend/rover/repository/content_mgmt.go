package repository

import (
	"context"

	"github.com/fidellr/jastip/backend/rover"
	"github.com/fidellr/jastip/backend/rover/models"
)

type ContentRepository interface {
	CreateScreenContent(ctx context.Context, m *models.Screen) error
	FetchContent(ctx context.Context, filter *rover.Filter) ([]*models.Screen, string, error)
	UpdateByContentID(ctx context.Context, shopID string, m *models.Screen) error
	GetContentByScreen(ctx context.Context, screenName string) (*models.Screen, error)
}
