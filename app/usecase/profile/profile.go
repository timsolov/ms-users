package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

// Repository describes repository contract
type Repository interface {
	// Profile returns user record
	Profile(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

// Profile describes parameters
type Profile struct {
	UserID uuid.UUID
}

// ProfileQuery describes dependencies
type ProfileQuery struct {
	repo Repository
}

func NewProfileQuery(repo Repository) ProfileQuery {
	return ProfileQuery{
		repo: repo,
	}
}

func (uc ProfileQuery) Do(ctx context.Context, query *Profile) (user entity.User, err error) {
	user, err = uc.repo.Profile(ctx, query.UserID)
	if err != nil {
		return user, err
	}

	return user, nil
}
