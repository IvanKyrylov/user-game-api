package game

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PointsGained int                `json:"points_gained,omitempty" bson:"points_gained,omitempty"`
	WinStatus    int8               `json:"win_status,omitempty" bson:"win_status,omitempty"`
	GameType     int8               `json:"game_type,omitempty" bson:"game_type,omitempty"`
	Created      time.Time          `json:"created,omitempty" bson:"created,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type GamesStatistics struct {
	GroupByDay []struct {
		GroupDate   string `json:"date" bson:"date"`
		GamesPlayed int64  `json:"games_played" bson:"games_played"`
	} `json:"group_by_day" bson:"group_by_day"`
	WithGameType []struct {
		GameDate    string `json:"date" bson:"date"`
		GameType    int8   `json:"game_type" bson:"game_type"`
		GamesPlayed int64  `json:"games_played" bson:"games_played"`
	} `json:"with_game_type" bson:"with_game_type"`
}
