package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserParams struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstname" bson:"firstname" binding:"required,min=2"`
	LastName  string             `json:"lastname" bson:"lastname" binding:"required,min=2"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Role      string             `json:"role" bson:"role"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	// ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func ToUserResponse(user *UserParams) UserResponse {
	return UserResponse{
		// ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// type UserServiceImpl struct {
// 	collection *mongo.Collection
// 	ctx        context.Context
// }

// func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
// 	return &UserServiceImpl{collection, ctx}
// }
