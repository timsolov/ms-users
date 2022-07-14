package postgres_test

import (
	"context"
	"ms-users/app/domain"
	"ms-users/app/repository/postgres"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/sjson"
)

func NewUser(t *testing.T, ctx context.Context, d *postgres.DB) (user domain.User, clean func()) {
	b, err := sjson.SetBytes(nil, "first_name", randomdata.FirstName(randomdata.Male))
	assert.NoError(t, err)
	b, err = sjson.SetBytes(b, "last_name", randomdata.LastName())
	assert.NoError(t, err)

	user = domain.User{
		UserID:  uuid.New(),
		View:    "test_fn_ln",
		Profile: b,
	}

	err = d.CreateUser(ctx, &user)
	assert.NoError(t, err)

	return user, func() {
		assert.NoError(t, d.DeleteUser(ctx, user.UserID))
	}
}
