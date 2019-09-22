package user

import (
	"context"
	"time"

	"github.com/fidellr/jastip/backend/uranus/models"
	"github.com/pkg/errors"

	uranus "github.com/fidellr/jastip/backend/uranus"
	"github.com/fidellr/jastip/backend/uranus/repository"
)

type service struct {
	repository     repository.UserAccountRepository
	validator      uranus.Validate
	contextTimeout time.Duration
}

func (s *service) CreateUserAccount(ctx context.Context, m *models.UserAccount) (err error) {
	if ctx == nil {
		err = uranus.ErrContextNil
		return err
	}

	if err = s.validator.ValidateStruct(m); err != nil {
		return err
	}

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	if err != nil {
		err = errors.Wrap(err, "error validating user account")
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err = s.repository.CreateUserAccount(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Fetch(ctx context.Context, filter *uranus.Filter) ([]*models.UserAccount, string, error) {
	if ctx == nil {
		err := uranus.ErrContextNil
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if filter.Num == 0 {
		filter.Num = int(20)
	}

	users, nextCursor, err := s.repository.Fetch(ctx, filter)
	if err != nil {
		return nil, nextCursor, err
	}

	return users, nextCursor, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*models.UserAccount, error) {
	if ctx == nil {
		err := uranus.ErrContextNil
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) SuspendAccount(ctx context.Context, id string) (bool, error) {
	if ctx == nil {
		err := uranus.ErrContextNil
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	isSuspsended, err := s.repository.SuspendAccount(ctx, id)
	if !isSuspsended || err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) RemoveAccount(ctx context.Context, id string) (bool, error) {
	if ctx == nil {
		err := uranus.ErrContextNil
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	isRemoved, err := s.repository.RemoveAccount(ctx, id)
	if !isRemoved || err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) UpdateUserByID(ctx context.Context, id string, m *models.UserAccount) (err error) {
	if ctx == nil {
		err = uranus.ErrContextNil
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	m.UpdatedAt = time.Now()

	err = s.repository.UpdateUserByID(ctx, id, m)
	if err != nil {
		return err
	}

	return nil
}

type requirement func(*service)

func Repository(repository repository.UserAccountRepository) requirement {
	return func(s *service) {
		s.repository = repository
	}
}

func Timeout(timeout time.Duration) requirement {
	return func(s *service) {
		s.contextTimeout = timeout
	}
}

func Validator(validator uranus.Validate) requirement {
	return func(s *service) {
		s.validator = validator
	}
}

func NewService(req ...requirement) uranus.UserAccountUsecase {
	s := new(service)
	for _, option := range req {
		option(s)
	}
	return s
}
