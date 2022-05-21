package postgres_test

import (
	"context"
	"ms-users/app/conf"
	"ms-users/app/domain"
	"ms-users/app/infrastructure/repository/postgres"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewIdent(t *testing.T, ctx context.Context, d *postgres.DB, userID uuid.UUID, kind domain.IdentKind) (ident domain.Ident, clean func()) {
	ident = domain.Ident{
		UserID:         userID,
		Kind:           kind,
		IdentConfirmed: true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	switch kind {
	case domain.EmailPassIdent:
		ident.Ident = randomdata.Email()
	default:
		t.Fatalf("unknown kind: %v", kind)
	}

	err := d.CreateIdent(ctx, &ident)
	assert.NoError(t, err)

	return ident, func() {
		assert.NoError(t, d.DeleteIdent(ctx, ident.Ident, ident.Kind))
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
