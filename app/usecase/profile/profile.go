package profile

import (
	"context"
	"ms-users/app/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Repository describes repository contract
type Repository interface {
	// User returns user record
	User(ctx context.Context, userID uuid.UUID) (domain.User, error)
	// IdentsByUserID returns idents for given user id.
	// If there is no idents for given user it returns empty list without error.
	IdentsByUserID(ctx context.Context, userID uuid.UUID) (idents []domain.Ident, err error)
}

// Params describes parameters
type Params struct {
	UserID uuid.UUID
}

// UseCase describes dependencies
type UseCase struct {
	repo Repository
}

func New(repo Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

func (uc UseCase) Do(ctx context.Context, query *Params) (domain.UserAggregate, error) {
	var ua domain.UserAggregate

	user, err := uc.repo.User(ctx, query.UserID)
	if err != nil {
		return ua, errors.Wrap(err, "request user record from db")
	}
	ua.User = user

	idents, err := uc.repo.IdentsByUserID(ctx, user.UserID)
	if err != nil {
		return ua, errors.Wrap(err, "request user's ident records from db")
	}

	ua.Idents = idents

	return ua, nil
}
