package repository

import (
	"context"
	"ms-users/app/domain/entity"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=../../infrastructure/repository/mockrepo/mockrepo.go -package=mockrepo ms-users/app/domain/repository Repository

// Repository is an API for work with database
type Repository interface {
	// CreateUserAggregate creates new ident record with user record.
	CreateUserAggregate(ctx context.Context, ua *entity.UserAggregate) error
	// Profile returns profile record
	Profile(ctx context.Context, profileID uuid.UUID) (entity.User, error)
	// EmailPassIdentByEmail returns email-pass identity by email.
	EmailPassIdentByEmail(ctx context.Context, email string) (ident entity.Ident, err error)
}
