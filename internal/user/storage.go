package user

import (
	"context"
)

type Storage interface {
	FindById(ctx context.Context, uuid string) (User, error)
	FindByName(ctx context.Context, lastName string) (User, error)
	FindAll(ctx context.Context, limit, page int64) ([]User, error)
}
