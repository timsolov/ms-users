package create_emailpass_identity

import (
	"context"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/common/password"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Repository describes repository contract
type Repository interface {
	// CreateUserAggregate creates new ident record with user record.
	CreateUserAggregate(ctx context.Context, ua *domain.UserAggregate, confirm *domain.Confirm, events []event.Event) error
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
	repo        Repository
	baseURL     string
	fromEmail   string
	fromName    string
	confirmLife time.Duration
}

func New(repo Repository, baseURL, fromEmail, fromName string, confirmLife time.Duration) UseCase {
	return UseCase{
		repo:        repo,
		baseURL:     baseURL,
		fromEmail:   fromEmail,
		fromName:    fromName,
		confirmLife: confirmLife,
	}
}

func (uc UseCase) Run(ctx context.Context, cmd *Params) (profileID uuid.UUID, err error) {
	// create encrypted password
	encryptedPass, err := password.Encrypt(cmd.Password)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "encrypting password by bcrypt")
	}

	// fill the profile
	profile := domain.V1Profile{
		Email:     cmd.Email,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
	}

	// TODO: when we'll have different identities we should check does identity
	//       already exist and attach new identity to existing user
	// new user record
	user := domain.User{
		UserID: uuid.New(),
		View:   "v1",
	}
	_ = user.MarshalProfile(profile)

	// new identity record
	ident := domain.Ident{
		UserID:         user.UserID,
		Ident:          cmd.Email,
		IdentConfirmed: false,
		Kind:           domain.EmailPassIdent,
		Password:       encryptedPass,
	}

	// fill user aggregate
	ua := domain.UserAggregate{
		User:   user,
		Idents: []domain.Ident{ident},
	}

	// prepare variables for email sending
	confirmRecord, err := domain.NewConfirm(
		domain.EmailConfirmKind,
		uc.confirmLife,
		nil, /* vars */
	)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "create new confirm struct")
	}
	confirmB64, err := confirmRecord.ToBase64()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "encode confirm struct to base64")
	}

	url := fmt.Sprintf("%s/confirm/%s", uc.baseURL, confirmB64)
	toEmail := cmd.Email
	toName := fmt.Sprintf("%s %s", cmd.FirstName, cmd.LastName)
	confirmEmail, err := event.EmailPassConfirm(
		"en", // TODO: en language should be user's language not constant
		uc.fromEmail,
		uc.fromName,
		toEmail,
		toName,
		url,
	)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "prepare email-pass confirm event")
	}

	err = uc.repo.CreateUserAggregate(ctx, &ua, &confirmRecord, []event.Event{confirmEmail})
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "create user record in db")
	}

	return user.UserID, nil
}
