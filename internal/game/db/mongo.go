package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IvanKyrylov/user-game-api/internal/apperror"
	"github.com/IvanKyrylov/user-game-api/internal/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	collection *mongo.Collection
	logger     *log.Logger
}

func NewStorage(storage *mongo.Database, collection string, logger *log.Logger) game.Storage {
	return &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
}

func (s *db) FindById(ctx context.Context, id string) (game game.Game, err error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return game, fmt.Errorf("failed to convert hex to objectid. error: %w", err)
	}

	filter := bson.M{"_id": objectId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := s.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		s.logger.Fatal(result.Err())
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return game, apperror.ErrNotFound
		}
		return game, fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err = result.Decode(&game); err != nil {
		return game, fmt.Errorf("failed to decode document. error: %w", err)
	}
	return game, nil
}

func (s *db) FindByPlayer(ctx context.Context, uuid string, limit, page int64) (games []game.Game, err error) {

	userId, err := primitive.ObjectIDFromHex(uuid)
	filter := bson.M{"user_uuid": userId}

	var opts *options.FindOptions
	if limit > 0 {
		opts.SetLimit(limit)
		opts.SetSkip(limit * page)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return games, apperror.ErrNotFound
		}
		return games, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = cur.All(ctx, &games); err == nil {
		return games, nil
	}
	return games, fmt.Errorf("failed to decode document. error: %w", err)
}

func (s *db) FindAll(ctx context.Context, limit, page int64) (games []game.Game, err error) {

	var opts *options.FindOptions
	if limit > 0 {
		opts.SetLimit(limit)
		opts.SetSkip(limit * page)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := s.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return games, apperror.ErrNotFound
		}
		return games, fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err = cur.All(ctx, &games); err == nil {
		return games, nil
	}
	return games, fmt.Errorf("failed to decode document. error: %w", err)
}
