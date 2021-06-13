package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UUID      primitive.ObjectID `json:"UUID,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Country   string             `json:"country,omitempty" bson:"country,omitempty"`
	City      string             `json:"city,omitempty" bson:"city,omitempty"`
	Gender    string             `json:"gender,omitempty" bson:"gender,omitempty"`
	BirthDate primitive.DateTime `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
}

type UserRating struct {
	User   User  `json:"user"`
	Rating int64 `json:"rating"`
}
