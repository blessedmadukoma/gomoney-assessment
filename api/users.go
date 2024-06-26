package api

import (
	"context"
	"strings"

	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (srv *Server) FindUserById(ctx context.Context, id string) (*db.UserParams, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *db.UserParams

	query := bson.M{"_id": oid}
	err := srv.Collections["users"].FindOne(ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &db.UserParams{}, err
		}
		return nil, err
	}

	return user, nil
}

func (srv *Server) FindUserByEmail(ctx context.Context, email string) (*db.UserParams, error) {
	var user db.UserParams

	query := bson.M{"email": strings.ToLower(email)}
	err := srv.Collections["users"].FindOne(ctx, query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
