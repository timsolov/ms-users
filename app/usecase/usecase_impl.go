package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/domain/repository"
)

// useCase describes usecases implementation.
type useCase struct {
	r repository.Repository
}

func New(r repository.Repository) Usecase {
	return &useCase{
		r: r,
	}
}

func (u *useCase) NewUser(ctx context.Context, user *entity.User) error {
	err := u.r.NewUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "create user in db")
	}

	return nil
}
