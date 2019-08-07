package repository

import (
	"context"

	"github.com/fidellr/jastip_way/backend/rover"
	"github.com/fidellr/jastip_way/backend/rover/models"
)

type ContentRepository interface {
	CreateScreenContent(ctx context.Context, m *models.Content) error
	FetchContent(ctx context.Context, filter *rover.Filter) ([]*models.Content, string, error)
	UpdateByContentID(ctx context.Context, shopID string, m *models.Content) error
	GetContentByScreen(ctx context.Context, screenName string) (*models.Content, error)
}
