package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

type UserRepository interface {
	// CreateUser creates new user record
	CreateUser(ctx context.Context, m *entity.User) error
}
