package seeds

import (
	"context"
	"log"
	"reflect"

	"github.com/blessedmadukoma/gomoney-assessment/api"
	"github.com/blessedmadukoma/gomoney-assessment/db"
	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	config      = utils.LoadEnvConfig(".env")
	collections = make(map[string]*mongo.Collection)
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})
)

var srv = *&api.Server{
	Config:      config,
	Collections: collections,
}

// Seed type
type Seed struct {
	collections map[string]*mongo.Collection
}

// Execute will executes the given seeder method
func Execute(config utils.Config, seedMethodNames ...string) {
	ctx := context.Background()

	_, collections := db.ConnectMongoDB(ctx, config)

	s := Seed{collections: collections}

	seedType := reflect.TypeOf(s)

	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		log.Println("Running all seeder...")
		for i := 0; i < seedType.NumMethod(); i++ {

			method := seedType.Method(i)

			seed(s, method.Name)
		}
	}

	// Execute only the given method names
	for _, item := range seedMethodNames {
		seed(s, item)
	}
}

func seed(s Seed, seedMethodName string) {

	m := reflect.ValueOf(s).MethodByName(seedMethodName)

	if !m.IsValid() {
		log.Fatal("No method called ", seedMethodName)
	}

	log.Println("Running", seedMethodName, "...")
	m.Call(nil)
	log.Println("Seed", seedMethodName, "was successful")
}
