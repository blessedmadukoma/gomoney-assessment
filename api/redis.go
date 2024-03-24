package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (srv *Server) GetDataFromRedis(ctx context.Context, key string, dest interface{}) (interface{}, error) {

	// val, err := srv.redisclient.Get(ctx, key).Result()
	val, err := srv.redisclient.Get(ctx, key).Bytes()
	if err != nil {
		return "", fmt.Errorf("failed to retreive my data from redis: %w", err)
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return "", fmt.Errorf("failed to unmarshal my data from JSON: %w", err)
	}

	return dest, nil
}

func (srv *Server) SetDataIntoRedis(ctx context.Context, key string, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data into JSON: %w", err)
	}

	// Set JSON value into Redis
	err = srv.redisclient.Set(ctx, key, jsonValue, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set data into Redis: %w", err)
	}
	return nil
}
