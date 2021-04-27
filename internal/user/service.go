package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
)

var _ Service = &service{}

type Service interface {
	GetById(ctx context.Context, uuid string) (User, error)
	GetByName(ctx context.Context, lastName string) (User, error)
	GetAll(ctx context.Context, limit, page int64) ([]User, error)
}

type service struct {
	storage Storage
	logger  *log.Logger
}

func NewService(storage Storage, logger *log.Logger) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
	}, nil
}

func (s service) GetById(ctx context.Context, uuid string) (user User, err error) {
	user, err = s.storage.FindById(ctx, uuid)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return user, err
		}
		return user, fmt.Errorf("failed to find user by uuid. error: %w", err)
	}
	return user, nil
}

func (s service) GetByName(ctx context.Context, lastName string) (user User, err error) {
	user, err = s.storage.FindByName(ctx, lastName)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return user, err
		}
		return user, fmt.Errorf("failed to get users by uuid. error: %w", err)
	}
	return user, nil
}

func (s service) GetAll(ctx context.Context, limit, page int64) (users []User, err error) {
	users, err = s.storage.FindAll(ctx, limit, page)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return users, err
		}
		return users, fmt.Errorf("failed to get all users. error: %w", err)
	}
	if len(users) == 0 {
		return users, apperror.ErrNotFound
	}
	return users, nil
}
