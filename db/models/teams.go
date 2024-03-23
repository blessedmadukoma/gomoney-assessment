package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamsParams struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TeamName  string             `json:"teamname" bson:"teamname" binding:"required,min=2"`
	ShortName string             `json:"shortname" bson:"shortname" binding:"required,min=2"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
