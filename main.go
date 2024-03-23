package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blessedmadukoma/gomoney-assessment/api"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ()

func main() {
	fmt.Println("Hello from GoMoney API! Starting Server...")

	config := utils.LoadEnvConfig(".env")
	// fmt.Println(config)

	ctx := context.Background()

	// connect to database
	mongoclient, redisclient := dbConn(ctx, config)

	value, err := redisclient.Get(ctx, "test").Result()
	defer mongoclient.Disconnect(ctx)

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	log.Println("value from redis:", value)

	// connect to database
	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }

	// store := db.NewStore(conn)
	// server, err := api.NewServer(config, store)

	// Get handles to the database and collections
	db := mongoclient.Database(config.MondoDBDatabase)

	collections := make(map[string]*mongo.Collection)

	// Add collections to the map
	collections["users"] = db.Collection("users")
	collections["teams"] = db.Collection("teams")
	collections["fixtures"] = db.Collection("fixtures")

	server, err := api.NewServer(config, collections)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// err = server.StartServer(config.ServerAddress)
	err = server.StartServer(config.Port)
	if err != nil {
		log.Fatal("cannot start server!")
	}
}

func dbConn(ctx context.Context, config utils.Config) (*mongo.Client, *redis.Client) {
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

	return mongoclient, redisclient

}
