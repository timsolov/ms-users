package postgres_test

import (
	"context"
	"fmt"
	"ms-users/app/conf"
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
		View:    randomdata.Noun(),
		Profile: b,
	}

	err = d.CreateUser(ctx, &user)
	assert.NoError(t, err)

	return user, func() {
		assert.NoError(t, d.DeleteUser(ctx, user.UserID))
	}
}

func TestUpdUser(t *testing.T) {
	config := conf.New()

	ctx := context.TODO()

	d, err := postgres.New(ctx, config.DB.DSN())
	assert.NoError(t, err)

	cases := []struct {
		name    string
		upd     domain.User
		columns []string
		wantErr error
		check   func(base, upd *domain.User)
	}{
		{
			name: "default_update_all",
			upd: domain.User{
				View:    "test",
				Profile: []byte(`{"first_name": "Name", "last_name": "Surname"}`),
			},
			wantErr: nil,
			check: func(base, upd *domain.User) {
				updated, err := d.User(ctx, base.UserID)
				if assert.NoError(t, err) {
					assert.Equal(t, upd.View, updated.View)
					assert.JSONEq(t, string(upd.Profile), string(updated.Profile))
				}
			},
		},
		{
			name: "only_view",
			upd: domain.User{
				View:    "test",
				Profile: []byte(`{"first_name": "Name", "last_name": "Surname"}`),
			},
			columns: []string{"view"},
			wantErr: nil,
			check: func(base, upd *domain.User) {
				updated, err := d.User(ctx, base.UserID)
				if assert.NoError(t, err) {
					assert.Equal(t, upd.View, updated.View)
					assert.JSONEq(t, string(base.Profile), string(updated.Profile))
				}
			},
		},
		{
			name: "only_profile",
			upd: domain.User{
				View:    "test",
				Profile: []byte(`{"first_name": "Name", "last_name": "Surname"}`),
			},
			columns: []string{"profile"},
			wantErr: nil,
			check: func(base, upd *domain.User) {
				updated, err := d.User(ctx, base.UserID)
				if assert.NoError(t, err) {
					assert.Equal(t, base.View, updated.View)
					assert.JSONEq(t, string(upd.Profile), string(updated.Profile))
				}
			},
		},
		{
			name:    "restricted",
			upd:     domain.User{},
			columns: []string{"created_at"},
			wantErr: fmt.Errorf("restricted columns for update: [created_at updated_at]"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user, clean := NewUser(t, ctx, d)
			defer clean()

			tc.upd.UserID = user.UserID

			err := d.UpdUser(ctx, &tc.upd, tc.columns...)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			}
			if tc.wantErr == nil && tc.check != nil {
				tc.check(&user, &tc.upd)
			}
		})
	}
}
