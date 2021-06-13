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
	if err != nil {
		return games, fmt.Errorf("failed to convert hex to objectid. error: %w", err)
	}
	filter := bson.M{"user_id": userId}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := s.collection.Find(ctx, filter, options.Find().SetLimit(limit).SetSkip(page*limit))
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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := s.collection.Find(ctx, bson.M{}, options.Find().SetLimit(limit).SetSkip(page*limit))
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

func (s *db) AggregateGamesStatistics(ctx context.Context, uuid string, startDate, endDate time.Time) (gamesStatistics []game.GamesStatistics, err error) {
	userId, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return gamesStatistics, fmt.Errorf("failed to convert hex to objectid. error: %w", err)
	}

	dateProjectStage := bson.M{
		"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created"},
	}

	groupedByDayStages := []bson.D{
		{{"$project", bson.M{
			"date":    dateProjectStage,
			"created": true,
		}}},
		{{"$group", bson.M{
			"_id":          "$date",
			"games_played": bson.M{"$sum": 1},
		}}},
		{{"$project", bson.M{
			"date":         "$_id",
			"games_played": true,
			"_id":          false,
		}}},
	}
	withGameTypeStages := []bson.D{
		{{"$project", bson.M{
			"date":      dateProjectStage,
			"created":   true,
			"game_type": true,
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"date":      "$date",
				"game_type": "$game_type",
			},
			"games_played": bson.M{"$sum": int64(1)},
		}}},
		{{"$project", bson.M{
			"date":         "$_id.date",
			"game_type":    "$_id.game_type",
			"games_played": true,
			"_id":          false,
		}}},
	}
	pipeline := []bson.D{
		{{"$match", bson.M{
			"user_id": userId,
			"created": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{"$facet", bson.M{
			"group_by_day":   groupedByDayStages,
			"with_game_type": withGameTypeStages,
		}}},
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return gamesStatistics, apperror.ErrNotFound
		}
		return gamesStatistics, fmt.Errorf("failed to execute query. error: %w", err)
	}

	// for cur.Next(ctx) {
	// 	var elem bson.M
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		return gamesStatistics, fmt.Errorf("failed to decode document. error: %w", err)
	// 	}

	// 	// groupByDay := elem["group_by_day"].(primitive.A)
	// 	// withGameType := elem["with_game_type"].(primitive.A)

	// 	gameStatistic := game.GamesStatistics{
	// 		GroupByDay: make([]struct {
	// 			GroupDate   time.Time "json:\"date\""
	// 			GamesPlayed int64     "json:\"games_played\""
	// 		}, len(elem["group_by_day"].(primitive.A))),
	// 		WithGameType: make([]struct {
	// 			GameDate    time.Time "json:\"date\""
	// 			GameType    int8      "json:\"game_type\""
	// 			GamesPlayed int64     "json:\"games_played\""
	// 		}, len(elem["with_game_type"].(primitive.A))),
	// 	}

	// 	for i, v := range elem["group_by_day"].(primitive.A) {
	// 		temp := v.()/(map[string]interface{})
	// 		gameStatistic.GroupByDay[i].GroupDate = temp["date"].(time.Time)
	// 		gameStatistic.GroupByDay[i].GamesPlayed = temp["games_played"].(int64)

	// 	}

	// 	for i, v := range elem["with_game_type"].(primitive.A) {
	// 		temp := v.(map[string]interface{})
	// 		gameStatistic.WithGameType[i].GameDate = temp["date"].(time.Time)
	// 		gameStatistic.WithGameType[i].GameType = temp["game_type"].(int8)
	// 		gameStatistic.WithGameType[i].GamesPlayed = temp["games_played"].(int64)
	// 	}

	// 	// temp := game.GamesStatistics{
	// 	// 	GroupByDay: struct {
	// 	// 		GroupDate   time.Time "json:\"date\""
	// 	// 		GamesPlayed int64     "json:\"games_played\""
	// 	// 	}{
	// 	// 		GroupDate:   groupByDay["date"].(time.Time),
	// 	// 		GamesPlayed: groupByDay["games_played"].(int64),
	// 	// 	},
	// 	// 	WithGameType: struct {
	// 	// 		GameDate    time.Time "json:\"date\""
	// 	// 		GameType    int8      "json:\"game_type\""
	// 	// 		GamesPlayed int64     "json:\"games_played\""
	// 	// 	}{
	// 	// 		GameDate:    withGameType["date"].(time.Time),
	// 	// 		GameType:    withGameType["game_type"].(int8),
	// 	// 		GamesPlayed: withGameType["games_played"].(int64),
	// 	// 	},
	// 	// }
	// 	gamesStatistics = gameStatistic

	// }
	// return gamesStatistics, fmt.Errorf("failed to decode document. error: %w", err)
	if err = cur.All(ctx, &gamesStatistics); err == nil {
		return gamesStatistics, nil
	}
	return gamesStatistics, fmt.Errorf("failed to decode document. error: %w", err)
}
