package testutils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/domain/repository"
)

type UserOption func(*entity.User)

func UserIDUserOption(userID uuid.UUID) UserOption {
	return func(user *entity.User) {
		user.UserID = userID
	}
}

// NewUser creates new user record in db
func NewUser(t *testing.T, r repository.Repository, opts ...UserOption) (user entity.User, clean func()) {
	ctx := context.TODO()

	user = entity.User{
		UserID: uuid.New(),
	}

	for _, opt := range opts {
		opt(&user)
	}

	assert.NoError(t, r.NewUser(ctx, &user))

	return user, func() {
		assert.NoError(t, r.DelUser(ctx, user.UserID))
	}
}
