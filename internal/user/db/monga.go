package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
	"github.com/IvanKyrylov/user-game-api/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewStorage(storage *mongo.Database, collection string, logger *log.Logger) user.Storage {
	return &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
}

func (s *db) FindById(ctx context.Context, uuid string) (user user.User, err error) {

	userId, err := primitive.ObjectIDFromHex(uuid)
	filter := bson.M{"_id": userId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result := s.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		s.logger.Println(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user, apperror.ErrNotFound
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode document. error: %w", err)
	}
	return user, nil

}

func (s *db) FindByName(ctx context.Context, lastName string) (user user.User, err error) {

	filter := bson.M{"last_name": lastName}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	result := s.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		s.logger.Println(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user, apperror.ErrNotFound
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = result.Decode(&user); err != nil {
		return user, fmt.Errorf("failed to decode document. error: %w", err)
	}
	return user, nil
}

func (s *db) FindAll(ctx context.Context, limit, page int64) (users []user.User, err error) {

	skip := page * limit
	opt := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	cur, err := s.collection.Find(ctx, bson.D{}, &opt)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return users, apperror.ErrNotFound
		}
		return users, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = cur.All(ctx, &users); err == nil {
		return users, nil
	}
	return users, fmt.Errorf("failed to decode document. error: %w", err)
}
