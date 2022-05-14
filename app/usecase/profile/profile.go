package profile

import (
	"context"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
)

// Repository describes repository contract
type Repository interface {
	// Profile returns user record
	Profile(ctx context.Context, userID uuid.UUID) (entity.User, error)
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

func (uc UseCase) Run(ctx context.Context, query *Params) (user entity.User, err error) {
	user, err = uc.repo.Profile(ctx, query.UserID)
	if err != nil {
		return user, err
	}

	return user, nil
}
