package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

type UserRepository interface {
	// NewUser creates new user record
	NewUser(ctx context.Context, m *entity.User) error
	// User returns user record by id
	User(ctx context.Context, userID uuid.UUID, columns ...string) (entity.User, error)
	// UpdUser changes user record's properties.
	UpdUser(ctx context.Context, m *entity.User, columns ...string) error
	// DelUser deletes user record record.
	DelUser(ctx context.Context, userID uuid.UUID) error
}
