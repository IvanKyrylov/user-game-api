package user

type User struct {
	UUID      string `json:"UUID,omitempty" bson:"_id,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
	LastName  string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Country   string `json:"country,omitempty" bson:"country,omitempty"`
	City      string `json:"city,omitempty" bson:"city,omitempty"`
	Gender    string `json:"gender,omitempty" bson:"gender,omitempty"`
	BirthDate string `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
	// GamesId   []string `json:"games_id,omitempty" bson:"games_id,omitempty"`
}
