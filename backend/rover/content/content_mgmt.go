package content

import (
	"context"
	"time"

	"github.com/fidellr/jastip/backend/rover"
	"github.com/fidellr/jastip/backend/rover/models"
	"github.com/fidellr/jastip/backend/rover/repository"
	"github.com/pkg/errors"
)

type service struct {
	repository     repository.ContentRepository
	validator      rover.Validate
	contextTimeout time.Duration
}

type requirement func(*service)

func Repository(repository repository.ContentRepository) requirement {
	return func(s *service) {
		s.repository = repository
	}
}

func Timeout(timeout time.Duration) requirement {
	return func(s *service) {
		s.contextTimeout = timeout
	}
}

func Validator(validator rover.Validate) requirement {
	return func(s *service) {
		s.validator = validator
	}
}

func NewService(reqs ...requirement) rover.ContentUsecase {
	s := new(service)
	for _, option := range reqs {
		option(s)
	}

	return s
}

func (s *service) CreateScreenContent(ctx context.Context, m *models.Content) (err error) {
	if ctx == nil {
		err = rover.ErrContextNil
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if err = s.validator.ValidateStruct(m); err != nil {
		return err
	}

	m.CreatedAt = time.Now()
	m.UpdateAt = time.Now()

	if err != nil {
		err = errors.Wrap(err, "error validating screen content")
		return err
	}

	err = s.repository.CreateScreenContent(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FetchContent(ctx context.Context, filter *rover.Filter) ([]*models.Content, string, error) {
	if ctx == nil {
		err := rover.ErrContextNil
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	if filter.Num == 0 {
		filter.Num = int(2)
	}

	content, nextCursor, err := s.repository.FetchContent(ctx, filter)
	if err != nil {
		return nil, nextCursor, err
	}

	return content, nextCursor, nil

}

func (s *service) GetContentByScreen(ctx context.Context, screenName string) (content *models.Content, err error) {
	if ctx == nil {
		err = rover.ErrContextNil
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	content, err = s.repository.GetContentByScreen(ctx, screenName)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (s *service) UpdateByContentID(ctx context.Context, screenID string, m *models.Content) (err error) {
	if ctx == nil {
		err = rover.ErrContextNil
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	m.UpdateAt = time.Now()

	err = s.repository.UpdateByContentID(ctx, screenID, m)
	if err != nil {
		return err
	}

	return nil
}
