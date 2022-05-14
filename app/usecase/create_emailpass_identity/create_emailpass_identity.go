package create_emailpass_identity

import (
	"context"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Repository describes repository contract
type Repository interface {
	// CreateUserAggregate creates new ident record with user record.
	CreateUserAggregate(ctx context.Context, ua *entity.UserAggregate) error
}

// Params describes parameters
type Params struct {
	Email          string
	EmailConfirmed bool
	FirstName      string
	LastName       string
	Password       string
}

// UseCase describes usecase
type UseCase struct {
	repo Repository
}

func New(repo Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

func (uc UseCase) Run(ctx context.Context, cmd *Params) (profileID uuid.UUID, err error) {
	encryptedPass, err := encrypt(cmd.Password)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "encrypting password by bcrypt")
	}

	profile := entity.V1Profile{
		Email:     cmd.Email,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
	}

	user := entity.User{
		UserID: uuid.New(),
		View:   "v1",
	}
	_ = user.MarshalProfile(profile)

	ident := entity.Ident{
		UserID:         user.UserID,
		Ident:          cmd.Email,
		IdentConfirmed: cmd.EmailConfirmed,
		Kind:           entity.EmailPassIdent,
		Password:       encryptedPass,
	}

	ua := entity.UserAggregate{
		User:   user,
		Idents: []entity.Ident{ident},
	}

	err = uc.repo.CreateUserAggregate(ctx, &ua)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "create user record in db")
	}

	// TODO: send email confirmation

	return user.UserID, nil
}

func encrypt(rawPass string) (encryptedPass string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.MinCost)
	return string(b), err
}
