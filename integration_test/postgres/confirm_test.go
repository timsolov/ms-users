package postgres_test

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/common/password"
	"ms-users/app/conf"
	"ms-users/app/domain"
	"ms-users/app/repository/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ConfirmOption func(c *domain.Confirm)

func EmailConfirmOption(email string) ConfirmOption {
	return func(c *domain.Confirm) {
		c.Vars["email"] = email
		c.Kind = domain.EmailConfirmKind
	}
}

func NewConfirm(t *testing.T, ctx context.Context, d *postgres.DB, opts ...ConfirmOption) (confirm domain.Confirm, clean func()) {
	plainPassword := uuid.NewString()
	encryptedPass, err := password.Encrypt(plainPassword)
	assert.NoError(t, err)

	confirm = domain.Confirm{
		ConfirmID:         uuid.New(),
		Password:          plainPassword,
		EncryptedPassword: encryptedPass,
		Kind:              domain.EmailConfirmKind,
		Vars: map[string]string{
			"email": "test@example.org",
		},
		CreatedAt: time.Now(),
		ValidTill: time.Now().Add(1 * time.Hour),
	}

	for _, opt := range opts {
		opt(&confirm)
	}

	err = d.CreateConfirm(ctx, &confirm, event.None)
	assert.NoError(t, err)

	return confirm, func() {
		assert.NoError(t, d.DelConfirm(ctx, confirm.ConfirmID))
	}
}

func Test_ReadConfirm(t *testing.T) {
	config := conf.New()

	ctx := context.TODO()

	d, err := postgres.New(ctx, config.DB.DSN())
	assert.NoError(t, err)

	confirm, clean := NewConfirm(t, ctx, d, EmailConfirmOption("test@example.org"))
	defer clean()

	_, err = d.ReadConfirm(ctx, confirm.ConfirmID)
	assert.NoError(t, err)
}
