package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{Client: rdb}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(key string) error {
	return r.Client.Del(ctx, key).Err()
}
