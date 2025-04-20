package cache

import (
	"context"
	"fmt"
	"socialAPI/internal/setting/cfg"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheStore interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg cfg.RedisConfig) (CacheStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверим подключение
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %v", err)
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), key, value, expiration).Err()
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *Redis) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *Redis) Exists(key string) (bool, error) {
	result, err := r.client.Exists(context.Background(), key).Result()
	return result > 0, err
}
