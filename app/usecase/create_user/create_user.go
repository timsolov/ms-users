package create_user

import (
	"context"
	"time"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
	Password  string
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
	encryptedPass, err := encrypt(cmd.Password)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "encrypting password by bcrypt")
	}

	user := entity.User{
		UserID:    uuid.New(),
		Email:     cmd.Email,
		Password:  encryptedPass,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.repo.CreateUser(ctx, &user)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "create user record in db")
	}

	return user.UserID, nil
}

func encrypt(rawPass string) (encryptedPass string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.MinCost)
	return string(b), err
}
