package repository

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../../infrastructure/repository/mockcache/mockcache.go -package=mockcache github.com/timsolov/ms-users/app/domain/repository Cache

type Cache interface {
	TTL(ctx context.Context, key string) (time.Duration, error)
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
}
