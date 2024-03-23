package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FixturesParams struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Home      string             `json:"home" bson:"home,omitempty" binding:"required,min=2"`
	Away      string             `json:"away" bson:"away,omitempty" binding:"required,min=2"`
	Status    string             `json:"status" bson:"status" binding:"min=2"`
	Fixture   string             `json:"fixture" bson:"fixture" binding:"min=2"`
	Link      string             `json:"link" bson:"link" binding:"min=2"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateFixturesParams struct {
	Home      string    `json:"home" bson:"home,omitempty" binding:"required,min=2"`
	Away      string    `json:"away" bson:"away,omitempty" binding:"required,min=2"`
	Status    string    `json:"status" bson:"status" binding:"min=2"`
	Fixture   string    `json:"fixture" bson:"fixture"`
	Link      string    `json:"link" bson:"link"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type FixturesResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Home      string             `json:"home" bson:"home"`
	Away      string             `json:"away" bson:"away"`
	Status    string             `json:"status" bson:"status"`
	Fixture   string             `json:"fixture" bson:"fixture"`
	Link      string             `json:"link" bson:"link"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func ToFixturesResponse(fixture *FixturesParams) FixturesResponse {
	return FixturesResponse{
		ID:        fixture.ID,
		Home:      fixture.Home,
		Away:      fixture.Away,
		Status:    fixture.Status,
		Fixture:   fixture.Fixture,
		Link:      fixture.Link,
		CreatedAt: fixture.CreatedAt,
		UpdatedAt: fixture.UpdatedAt,
	}
}
