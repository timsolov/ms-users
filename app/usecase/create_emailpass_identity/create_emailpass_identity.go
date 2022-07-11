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

const (
	confirmEmailTemplateName = "confirm-registration-email-pass"
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

func (uc UseCase) Do(ctx context.Context, cmd *Params) (profileID uuid.UUID, err error) {
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

	confirmRecord, err := uc.PrepareConfirmEmailRecord(cmd.Email)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "prepare confirm email record") // 500
	}

	// encode to json object which encoded by base64
	code, err := confirmRecord.ToBase64JSON()
	if err != nil {
		err = errors.Wrap(err, "encode confirm struct to base64") // 500
		return
	}

	confirmEmailEvent, err := uc.PrepareConfirmEmailEvent(cmd.FirstName, cmd.LastName, cmd.Email, "en", code)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "prepare confirm email event") // 500
	}

	err = uc.repo.CreateUserAggregate(ctx, &ua, &confirmRecord, []event.Event{confirmEmailEvent})
	if err != nil {
		if errors.Cause(err) == domain.ErrNotUnique {
			err = domain.ErrIdentityDuplicated // 400
		}
		return uuid.Nil, errors.Wrap(err, "create user record in db") // 500
	}

	return user.UserID, nil
}

// PrepareConfirmEmailRecord creates instance of Confirm struct for confirming email.
func (uc UseCase) PrepareConfirmEmailRecord(email string) (confirm domain.Confirm, err error) {
	// prepare variables for email sending
	confirmPassword := uuid.New().String()
	confirm, err = domain.NewConfirm(
		domain.EmailConfirmKind,
		confirmPassword,
		uc.confirmLife,
		map[string]string{ /* vars */
			"email": email,
		},
	)

	return
}

// PrepareConfirmEmailEvent creates instance of Event for sending confirmation email.
func (uc UseCase) PrepareConfirmEmailEvent(firstName, lastName, email, lang, code string) (confirmEmail event.Event, err error) {
	// prepare url
	url := fmt.Sprintf("%s/confirm/%s", uc.baseURL, code)

	// prepare event email.SendTemplate
	toEmail := email
	toName := fmt.Sprintf("%s %s", firstName, lastName)
	confirmEmail, err = event.EmailSendTemplate(
		confirmEmailTemplateName,
		lang, // TODO: en language should be user's language not constant
		uc.fromEmail,
		uc.fromName,
		toEmail,
		toName,
		map[string]string{
			"url": url,
		},
	)
	if err != nil {
		err = errors.Wrap(err, "prepare event for sending email-pass confirmation")
		return
	}

	return confirmEmail, err
}
