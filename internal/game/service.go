package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
)

var _ Service = &service{}

type Service interface {
	GetById(ctx context.Context, id string) (Game, error)
	GetByPlayer(ctx context.Context, uuid string, limit, page int64) ([]Game, error)
	GetAll(ctx context.Context, limit, page int64) ([]Game, error)
	GetGamesStatistics(ctx context.Context, userId string, startDate, endDate time.Time) ([]GamesStatistics, error)
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

func (s service) GetById(ctx context.Context, id string) (game Game, err error) {
	game, err = s.storage.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound.Err) {
			return game, err
		}
		return game, fmt.Errorf("failed to find game by id. error: %w", err)
	}
	return game, nil
}

func (s service) GetByPlayer(ctx context.Context, uuid string, limit, page int64) (games []Game, err error) {
	games, err = s.storage.FindByPlayer(ctx, uuid, limit, page)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return games, err
		}
		return games, fmt.Errorf("failed to get games by player. error: %w", err)
	}
	if len(games) == 0 {
		return games, apperror.ErrNotFound
	}
	return games, nil
}

func (s service) GetAll(ctx context.Context, limit, page int64) (games []Game, err error) {
	games, err = s.storage.FindAll(ctx, limit, page)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return games, err
		}
		return games, fmt.Errorf("failed to get all games. error: %w", err)
	}
	if len(games) == 0 {
		return games, apperror.ErrNotFound
	}
	return games, nil
}

func (s service) GetGamesStatistics(ctx context.Context, userId string, startDate, endDate time.Time) (data []GamesStatistics, err error) {

	data, err = s.storage.AggregateGamesStatistics(ctx, userId, startDate, endDate)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return data, err
		}
		return data, fmt.Errorf("failed to get statistics games. error: %w", err)
	}
	if len(data) == 0 {
		return data, apperror.ErrNotFound
	}

	return data, nil
}
