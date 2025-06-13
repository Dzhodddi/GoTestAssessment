package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type CatStore struct {
	cacheRedis *redis.Client
}

func (s *CatStore) Get(ctx context.Context, breed string) (bool, error) {
	data, err := s.cacheRedis.Get(ctx, breed).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return strconv.ParseBool(data)
}

func (s *CatStore) Set(ctx context.Context, breed string) error {
	return s.cacheRedis.Set(ctx, breed, true, time.Hour).Err()
}
