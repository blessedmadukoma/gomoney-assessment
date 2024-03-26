package api

import (
	"context"
	"encoding/json"
	"time"
)

func (srv *Server) GetDataFromRedis(ctx context.Context, key string, dest interface{}) (interface{}, error) {

	val, err := srv.redisclient.Get(ctx, key).Bytes()
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(val, dest); err != nil {
		return "", err
	}

	return dest, nil
}

func (srv *Server) SetDataIntoRedis(ctx context.Context, key string, value any) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set JSON value into Redis
	err = srv.redisclient.Set(ctx, key, jsonValue, 3*time.Minute).Err() // expire eery 3 minutes
	if err != nil {
		return err
	}
	return nil
}
