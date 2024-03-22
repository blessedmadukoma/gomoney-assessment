package db_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestMain(m *testing.M) {

	config := utils.LoadEnvConfig("../../.env")

	ctx := context.Background()

	log.Println("Config in testdb:", config)

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.MongoDBSource)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// ? Connect to Redis
	redisclient := redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB",
		0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")
}
