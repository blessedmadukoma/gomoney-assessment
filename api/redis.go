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
		return "", fmt.Errorf("failed to retreive data from redis: %w", err)
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return "", fmt.Errorf("failed to unmarshal data from JSON: %w", err)
	}

	return dest, nil
}

func (srv *Server) SetDataIntoRedis(ctx context.Context, key string, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set JSON value into Redis
	err = srv.redisclient.Set(ctx, key, jsonValue, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
