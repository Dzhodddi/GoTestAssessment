package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Cats interface {
		Get(context.Context, string) (bool, error)
		Set(context.Context, string) error
	}
}

func NewRedisStorage(cacheRedis *redis.Client) Storage {
	return Storage{
		Cats: &CatStore{cacheRedis: cacheRedis},
	}
}
