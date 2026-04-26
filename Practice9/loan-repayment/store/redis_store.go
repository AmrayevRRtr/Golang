package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisStore{client: rdb}
}

func (r *RedisStore) Start(ctx context.Context, key string) (bool, error) {
	return r.client.SetNX(ctx, key, "processing", 5*time.Minute).Result()
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStore) Save(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 24*time.Hour).Err()
}
