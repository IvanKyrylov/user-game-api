package game

type Game struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	PointsGained string `json:"points_gained,omitempty" bson:"points_gained,omitempty"`
	WinStatus    string `json:"win_status,omitempty" bson:"win_status,omitempty"`
	GameType     string `json:"game_type,omitempty" bson:"game_type,omitempty"`
	Created      string `json:"created,omitempty" bson:"created,omitempty"`
	// UserUUID     string `json:"user_uuid,omitempty" bson:"user_uuid,omitempty"`
}
