package usecase

import (
	"context"

	"github.com/timsolov/ms-users/app/domain/entity"
)

type Usecase interface {
	// NewUser creates new user record
	NewUser(ctx context.Context, user *entity.User) error
}
