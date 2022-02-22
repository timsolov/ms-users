package testutils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timsolov/ms-users/internal/client/db"
	"github.com/timsolov/ms-users/internal/entity"
)

type UserOption func(*entity.User)

func UserIDUserOption(userID uuid.UUID) UserOption {
	return func(user *entity.User) {
		user.UserID = userID
	}
}

// NewUser creates new user record in db
func NewUser(t *testing.T, d db.DB, opts ...UserOption) (entity.User, func()) {
	ctx := context.TODO()

	user := entity.User{
		UserID: uuid.New(),
	}

	for _, opt := range opts {
		opt(&user)
	}

	assert.NoError(t, d.NewUser(ctx, &user))

	return user, func() {
		assert.NoError(t, d.DelUser(ctx, user.UserID))
	}
}
