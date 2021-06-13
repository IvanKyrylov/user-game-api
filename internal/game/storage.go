package game

import (
	"context"
	"time"
)

type Storage interface {
	FindById(ctx context.Context, id string) (Game, error)
	FindByPlayer(ctx context.Context, uuid string, limit, page int64) ([]Game, error)
	FindAll(ctx context.Context, limit, page int64) ([]Game, error)
	AggregateGamesStatistics(ctx context.Context, uuid string, startDate, endDate time.Time) ([]GamesStatistics, error)
}
