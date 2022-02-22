package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/timsolov/ms-users/internal/client/db"
	"github.com/timsolov/ms-users/internal/entity"
)

// stub: u *useCase usecase.UseCase

type UseCase interface {
	entity.UserUseCase
}

// useCase describes usecases implementation.
type useCase struct {
	db db.DB
}

func New(db db.DB) UseCase {
	return &useCase{
		db: db,
	}
}

func (u *useCase) NewUser(ctx context.Context, user *entity.User) error {
	err := u.db.NewUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "create user in db")
	}

	return nil
}
