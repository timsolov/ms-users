package repository

import (
	"context"
)

//go:generate mockgen -destination=../../infrastructure/repository/mockrepo/mockrepo.go -package=mockrepo github.com/timsolov/ms-users/app/domain/repository Repository

// Repository is an API for work with database
type Repository interface {
	// HealthCheck returns database health check.
	HealthCheck() error
	// Atomic executes operations in transaction.
	Atomic(ctx context.Context, fn func(tx Repository) error) error

	UserRepository
}
