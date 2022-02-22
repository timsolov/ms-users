package db

import (
	"context"

	"github.com/timsolov/ms-users/internal/entity"
)

//-go:generate mockgen -source=db.go -destination=./testdb/mock_db.go -package=testdb
//-go:generate sqlc generate

// DB is an API for work with database
type DB interface {
	// HealthCheck returns database health check.
	HealthCheck() error
	// Atomic executes operations in transaction.
	Atomic(ctx context.Context, fn func(tx DB) error) error

	entity.UserRepository
}
