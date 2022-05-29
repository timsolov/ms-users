package postgres_test

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/conf"
	"ms-users/app/domain"
	"ms-users/app/repository/postgres"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func NewIdent(t *testing.T, ctx context.Context, d *postgres.DB, userID uuid.UUID, kind domain.IdentKind) (ident domain.Ident, clean func()) {
	ident = domain.Ident{
		UserID:    userID,
		Kind:      kind,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	switch kind {
	case domain.EmailPassIdent:
		ident.Ident = randomdata.Email()
	default:
		t.Fatalf("unknown kind: %v", kind)
	}

	err := d.CreateIdent(ctx, &ident)
	assert.NoError(t, err)

	i, k := ident.Ident, ident.Kind

	return ident, func() {
		assert.NoError(t, d.DeleteIdent(ctx, i, k))
	}
}

func TestDB_EmailPassIdentByEmail(t *testing.T) {
	config := conf.New()

	ctx := context.TODO()

	d, err := postgres.New(ctx, config.DB.DSN())
	assert.NoError(t, err)

	user, clean := NewUser(t, ctx, d)
	defer clean()

	ident, clean := NewIdent(t, ctx, d, user.UserID, domain.EmailPassIdent)
	defer clean()

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx:   ctx,
				email: ident.Ident,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := d.EmailPassIdentByEmail(tt.args.ctx, tt.args.email)
			if err != tt.wantErr {
				t.Errorf("DB.EmailPassIdentByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDB_UpdIdent(t *testing.T) {
	config := conf.New()

	ctx := context.TODO()

	d, err := postgres.New(ctx, config.DB.DSN())
	assert.NoError(t, err)

	user, clean := NewUser(t, ctx, d)
	defer clean()

	tests := []struct {
		name    string
		ident   string
		kind    domain.IdentKind
		create  bool
		update  func(ident *domain.Ident)
		check   func(ident *domain.Ident)
		wantErr error
	}{
		{
			name:   "ident_confirmed_true",
			ident:  randomdata.Email(),
			kind:   domain.EmailPassIdent,
			create: true,
			update: func(ident *domain.Ident) {
				ident.IdentConfirmed = true
			},
			check: func(ident *domain.Ident) {
				assert.Equal(t, true, ident.IdentConfirmed)
			},
			wantErr: nil,
		},
		{
			name:   "ident_not_found",
			ident:  randomdata.Email(),
			kind:   domain.EmailPassIdent,
			create: false,
			update: func(ident *domain.Ident) {
				ident.Ident = randomdata.Email()
			},
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				ident domain.Ident
				clean func()
			)

			if tt.create {
				ident, clean = NewIdent(t, ctx, d, user.UserID, tt.kind)
				defer clean()
			}

			if tt.update != nil {
				tt.update(&ident)
			}

			err := d.UpdIdent(ctx, &ident, event.None)

			if assert.Equal(t, tt.wantErr, errors.Cause(err)) {
				if err == nil {
					ik, err := d.ReadIdentKind(ctx, ident.Ident, ident.Kind)
					assert.NoError(t, err)

					if tt.check != nil {
						tt.check(&ik)
					}
				}
			} else {
				t.Errorf("DB.EmailPassIdentByEmail() error actual = %v, wantErr = %v", errors.Cause(err), tt.wantErr)
				return
			}
		})
	}
}
