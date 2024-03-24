package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blessedmadukoma/gomoney-assessment/api"
	"github.com/blessedmadukoma/gomoney-assessment/db"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

var ()

func main() {
	fmt.Println("Hello from GoMoney API! Starting Server...")

	config := utils.LoadEnvConfig(".env")

	ctx := context.Background()

	// connect to database
	mongoclient, redisclient := dbConn(ctx, config)

	// value, err := redisclient.Get(ctx, "teams").Result()
	value, err := redisclient.Get(ctx, "test").Result()
	defer mongoclient.Disconnect(ctx)

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		log.Fatal("unable to get value from redis client:", err)
		return
	}

	log.Println("value from redis:", value)

	// Get handles to the database and collections
	db := mongoclient.Database(config.MondoDBDatabase)

	collections := make(map[string]*mongo.Collection)

	// Add collections to the map
	collections["users"] = db.Collection("users")
	collections["teams"] = db.Collection("teams")
	collections["fixtures"] = db.Collection("fixtures")

	server, err := api.NewServer(config, collections, redisclient)
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
	// mongoconn := options.Client().ApplyURI(config.MongoDBSource)
	// mongoclient, err := mongo.Connect(ctx, mongoconn)

	// if err != nil {
	// 	log.Fatal("unable to connect to mongodb:", err)
	// 	return nil, nil
	// }

	// if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
	// 	log.Fatal("unable to ping mongodb:", err)
	// 	return nil, nil
	// }

	// fmt.Println("MongoDB successfully connected...")

	mongoclient, _ := db.ConnectMongoDB(ctx, config)

	redisclient := db.ConnectRedis(ctx, config)

	return mongoclient, redisclient

}
