package db

import (
	"context"
	"fmt"
	"log"

	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB(ctx context.Context, config utils.Config) (*mongo.Client, map[string]*mongo.Collection) {
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

func ConnectRedis(ctx context.Context, config utils.Config) *redis.Client {
	// ? Connect to Redis
	redisclient := redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		log.Fatal("unable to ping redis:", err)
		return nil
	}

	err := redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB",
		0).Err()
	if err != nil {
		log.Fatal("unable to set value in redis:", err)
		return nil
	}

	fmt.Println("Redis client connected successfully...")

	return redisclient
}
