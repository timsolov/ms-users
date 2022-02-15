package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timsolov/ms-users/app/conf"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/infrastructure/repository/postgres"
	"github.com/timsolov/ms-users/app/infrastructure/repository/postgres/testutils"
)

func TestDB_NewUser_DelUser(t *testing.T) {
	cfg := conf.New()

	ctx := context.Background()

	d, err := postgres.New(ctx, cfg.DB().DSN(), 1, 1, 1*time.Minute)
	assert.NoError(t, err)

	user := entity.User{
		UserID: uuid.New(),
	}
	err = d.NewUser(ctx, &user)
	if assert.NoError(t, err) {
		assert.NoError(t, d.DelUser(ctx, user.UserID))
	}
}

func TestDB_User(t *testing.T) {
	cfg := conf.New()

	ctx := context.Background()

	d, err := postgres.New(ctx, cfg.DB().DSN(), 1, 1, 1*time.Minute)
	assert.NoError(t, err)

	user, clean := testutils.NewUser(t, d)
	defer clean()

	newUser, err := d.User(ctx, user.UserID)
	if assert.NoError(t, err) {
		assert.Equal(t, user.UserID, newUser.UserID)
	}
}

func TestDB_UpdUser(t *testing.T) {
	cfg := conf.New()

	ctx := context.Background()

	d, err := postgres.New(ctx, cfg.DB().DSN(), 1, 1, 1*time.Minute)
	assert.NoError(t, err)

	user, clean := testutils.NewUser(t, d)
	defer clean()

	user.FirstName = randomdata.Alphanumeric(10)

	err = d.UpdUser(ctx, &user, "name")
	if assert.NoError(t, err) {
		newUser, err := d.User(ctx, user.UserID)
		if assert.NoError(t, err) {
			assert.Equal(t, user.UserID, newUser.UserID)
			assert.Equal(t, user.FirstName, newUser.FirstName)
		}
	}
}
