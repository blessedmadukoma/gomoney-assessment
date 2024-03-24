package api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	"github.com/blessedmadukoma/gomoney-assessment/utils"
	"github.com/go-redis/redis/v8"
)

var config = utils.LoadEnvConfig("../.env")

func TestGetDataFromRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})

	srv := &Server{
		redisclient: redisClient,
	}

	testData := map[string]interface{}{
		"test": "get data from redis",
	}

	jsonData, _ := json.Marshal(testData)
	err := redisClient.Set(context.Background(), "testKey", jsonData, time.Minute).Err()
	assert.NoError(t, err)

	var result map[string]interface{}
	_, err = srv.GetDataFromRedis(context.Background(), "testKey", &result)
	assert.NoError(t, err)
	assert.Equal(t, testData, result)
}

func TestSetDataIntoRedis(t *testing.T) {
	// Initialize a real Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisDBSource,
	})

	srv := &Server{
		redisclient: redisClient,
	}

	testData := map[string]interface{}{
		"test": "set data into redis",
	}

	err := srv.SetDataIntoRedis(context.Background(), "testKey", testData)
	assert.NoError(t, err)

	val, err := redisClient.Get(context.Background(), "testKey").Result()
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(val), &result)
	assert.NoError(t, err)
	assert.Equal(t, testData, result)
}
