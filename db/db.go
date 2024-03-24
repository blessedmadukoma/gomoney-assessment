package db

import (
	"context"
	"fmt"
	"log"

	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoConn(ctx context.Context, config utils.Config) (*mongo.Client, map[string]*mongo.Collection) {
	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.MongoDBSource)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Fatal("unable to connect to mongodb:", err)
		return nil, nil
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("unable to ping mongodb:", err)
		return nil, nil
	}

	fmt.Println("MongoDB successfully connected...")

	db := mongoclient.Database(config.MondoDBDatabase)

	collections := make(map[string]*mongo.Collection)

	// Add collections to the map
	collections["users"] = db.Collection("users")
	collections["teams"] = db.Collection("teams")
	collections["fixtures"] = db.Collection("fixtures")

	return mongoclient, collections
}
