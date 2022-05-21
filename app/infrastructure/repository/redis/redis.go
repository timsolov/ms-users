package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-users/app/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func New(redisDSN string) (*Redis, error) {
	opt, err := redis.ParseURL(redisDSN)
	if err != nil {
		return nil, fmt.Errorf("parse redis DSN: %w", err)
	}

	redisStruct := &Redis{
		client: redis.NewClient(opt),
	}

	ctx := context.Background()
	if err := redisStruct.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return redisStruct, nil
}

// TTL returns the remaining time to live of a key that has a timeout.
// This introspection capability allows a Redis client to check how many seconds a given key will continue to be part of the dataset.
func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	if ttl < 0 {
		return -1, domain.ErrNotFound
	}

	return ttl, nil
}

// Get retrieves value by given key from Redis
func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	// get string from redis
	s, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.ErrNotFound
		}
		return err
	}

	// deserialize string into
	err = json.Unmarshal([]byte(s), value)
	if err != nil {
		return err
	}

	return nil
}

// Set sets value to specified key
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// serialize value object
	s, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// save to redis
	if err := r.client.Set(ctx, key, s, expiration).Err(); err != nil {
		return err
	}

	return nil
}

// Del deletes value by key
func (r *Redis) Del(ctx context.Context, key string) error {
	// delete from redis
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

// Ping sends keep-alive message
func (r *Redis) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("ping result: %w", err)
	}

	return nil
}
