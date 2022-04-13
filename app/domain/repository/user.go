package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

type UserRepository interface {
	// CreateUser creates new user record
	CreateUser(ctx context.Context, m *entity.User) error
	// Profile returns user record
	Profile(ctx context.Context, userID uuid.UUID) (entity.User, error)
}
