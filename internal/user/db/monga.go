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

	cur, err := s.collection.Find(ctx, bson.M{}, options.Find().SetLimit(limit).SetSkip(page*limit))

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

func (s *db) AggregateRatingUsers(ctx context.Context, limit, page int64) (usersRatings []user.UserRating, err error) {
	skip := page * limit

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// pipeline := make([]bson.M, 0)
	// groupStage := bson.M{
	// 	"$group": bson.M{
	// 		"_id":        "$user_id",
	// 		"count_game": bson.M{"$sum": 1},
	// 	},
	// }

	// pipeline = append(pipeline, groupStage)

	// pipeline := []bson.M{
	// 	{
	// 		"$group": bson.M{
	// 			"_id":         "$user_id",
	// 			"count_games": bson.M{"$sum": 1},
	// 		},
	// 	},
	// 	{
	// 		"$limit": limit,
	// 	},
	// 	{
	// 		"$skip": skip,
	// 	},
	// }

	// id, _ := primitive.ObjectIDFromHex("60c0928d9bfeb2397e77124d")
	// matchStage := bson.D{{"$match", bson.D{{"user_id", id}}}}

	// skipStage := bson.D{{"$skip", skip}}
	// limitStage := bson.D{{"$limit", limit}}
	// groupStage := bson.D{{"$group", bson.D{{"_id", "$user_id"}, {"count_games", bson.D{{"$sum", 1}}}}}}
	// sortStage := bson.D{{"$sort", bson.D{{"count_game", -1}}}}

	cur, err := s.collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{"rating", -1}}).SetLimit(limit).SetSkip(skip))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return usersRatings, apperror.ErrNotFound
		}
		return usersRatings, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for cur.Next(ctx) {
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			return usersRatings, fmt.Errorf("failed to decode document. error: %w", err)
		}

		userField := user.User{
			UUID:      elem["_id"].(primitive.ObjectID),
			Email:     elem["email"].(string),
			LastName:  elem["last_name"].(string),
			Country:   elem["country"].(string),
			City:      elem["city"].(string),
			Gender:    elem["gender"].(string),
			BirthDate: elem["birth_date"].(primitive.DateTime),
		}

		userRating := user.UserRating{
			User:   userField,
			Rating: elem["rating"].(int64),
		}
		usersRatings = append(usersRatings, userRating)
	}

	return usersRatings, nil

}
