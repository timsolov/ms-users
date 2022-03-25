package create_user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/domain/repository"
)

// CreateUser describes
type CreateUser struct {
	Email     string
	FirstName string
	LastName  string
}

// CreateUserCommand describes
type CreateUserCommand struct {
	repo repository.UserRepository
}

func NewCreateUserCommand(repo repository.UserRepository) CreateUserCommand {
	return CreateUserCommand{
		repo: repo,
	}
}

func (uc CreateUserCommand) Do(ctx context.Context, cmd *CreateUser) (userID uuid.UUID, err error) {
	user := entity.User{
		UserID:    uuid.New(),
		Email:     cmd.Email,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.repo.NewUser(ctx, &user)
	if err != nil {
		return uuid.Nil, err
	}

	return user.UserID, nil
}
