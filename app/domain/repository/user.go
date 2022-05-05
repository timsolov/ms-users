package repository

import (
	"context"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	// CreateUser creates new user record
	CreateUser(ctx context.Context, m *entity.User) error
	// Profile returns user record
	Profile(ctx context.Context, userID uuid.UUID) (entity.User, error)
}
