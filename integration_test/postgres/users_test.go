package postgres_test

import (
	"context"
	"encoding/json"
	"ms-users/app/domain/entity"
	"ms-users/app/infrastructure/repository/postgres"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewUser(t *testing.T, ctx context.Context, d *postgres.DB) (user entity.User, clean func()) {
	p := entity.V1Profile{
		Email:     randomdata.Email(),
		FirstName: randomdata.FirstName(randomdata.Male),
		LastName:  randomdata.LastName(),
	}

	b, err := json.Marshal(p)
	assert.NoError(t, err)

	user = entity.User{
		UserID:  uuid.New(),
		View:    "v1",
		Profile: b,
	}

	err = d.CreateUser(ctx, &user)
	assert.NoError(t, err)

	return user, func() {
		assert.NoError(t, d.DeleteUser(ctx, user.UserID))
	}
}
