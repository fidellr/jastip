package repository

import (
	"context"

	"github.com/fidellr/jastip/backend/uranus"

	"github.com/fidellr/jastip/backend/uranus/models"
)

// UserAccountRepository repo
type UserAccountRepository interface {
	CreateUserAccount(ctx context.Context, userAccountM *models.UserAccount) error
	Fetch(ctx context.Context, filter *uranus.Filter) ([]*models.UserAccount, string, error)
	GetUserByID(ctx context.Context, uuid string) (*models.UserAccount, error)
	SuspendAccount(ctx context.Context, uuid string) (bool, error)
	RemoveAccount(ctx context.Context, uuid string) (bool, error)
	UpdateUserByID(ctx context.Context, uuid string, userAccountM *models.UserAccount) error
}
