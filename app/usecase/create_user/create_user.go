package create_user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

// Repository describes repository contract
type Repository interface {
	// CreateUser creates new user record
	CreateUser(ctx context.Context, m *entity.User) error
}

// CreateUser describes parameters
type CreateUser struct {
	Email     string
	FirstName string
	LastName  string
}

// CreateUserCommand describes dependencies
type CreateUserCommand struct {
	repo Repository
}

func NewCreateUserCommand(repo Repository) CreateUserCommand {
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

	err = uc.repo.CreateUser(ctx, &user)
	if err != nil {
		return uuid.Nil, err
	}

	return user.UserID, nil
}
