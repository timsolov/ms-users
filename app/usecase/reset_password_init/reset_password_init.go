package reset_password_init

import (
	"context"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/common/utils"
	"ms-users/app/domain"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const (
	resetEmailPasswordTpl = "reset-email-password-init"
)

var (
	validation = validator.New()
)

// Repository describes repository contract
type Repository interface {
	// ReadIdentKind returns ident by ident and kind.
	ReadIdentKind(ctx context.Context, ident string, kind domain.IdentKind) (domain.Ident, error)
	// Profile returns profile record
	Profile(ctx context.Context, userID uuid.UUID) (domain.User, error)
	// CreateConfirm creates new confirm record
	CreateConfirm(ctx context.Context, m *domain.Confirm, events []event.Event) error
}

// Params describes parameters
type Params struct {
	Ident string
}

// UseCase describes usecase
type UseCase struct {
	repo        Repository
	webBaseURL  string
	fromEmail   string
	fromName    string
	confirmLife time.Duration
}

// New creates new instance of the UseCase
func New(repo Repository, webBaseURL, fromEmail, fromName string, confirmLife time.Duration) UseCase {
	return UseCase{
		repo:        repo,
		webBaseURL:  webBaseURL,
		fromEmail:   fromEmail,
		fromName:    fromName,
		confirmLife: confirmLife,
	}
}

// Do the main method for the UseCase
func (uc UseCase) Do(ctx context.Context, cmd *Params) (err error) {
	const (
		isEmail = "required,email"
	)
	kind := domain.UnknownIdent
	if validation.VarCtx(ctx, cmd.Ident, isEmail) == nil {
		kind = domain.EmailPassIdent
	}

	// searching for the ident
	var ident domain.Ident

	switch kind {
	case domain.EmailPassIdent:
		ident, err = uc.repo.ReadIdentKind(ctx, cmd.Ident, domain.EmailPassIdent)
	default:
		return domain.ErrUnknownIdent
	}
	if err != nil {
		if errors.Cause(err) == domain.ErrNotFound {
			err = domain.ErrEmailPassNotFound // 204
		}
		err = errors.Wrap(err, "request email-pass ident from db") // 500
		return
	}

	// searching for the user profile
	profile, err := uc.repo.Profile(ctx, ident.UserID)
	if err != nil {
		return errors.Wrap(err, "request user profile from db") // 500
	}

	firstName := gjson.GetBytes(profile.Profile, "first_name").String()
	lastName := gjson.GetBytes(profile.Profile, "last_name").String()
	if firstName == "" || lastName == "" {
		return fmt.Errorf("first_name or last_name is empty") // 500
	}

	// prepare confirm record and email event
	var (
		confirmRecord domain.Confirm
		confirmEmail  event.Event
	)
	if kind == domain.EmailPassIdent {
		confirmRecord, confirmEmail, err = uc.prepareConfirmRecordWithEvent(ctx, cmd.Ident, firstName, lastName, "en")
	}
	if err != nil {
		return errors.Wrap(err, "prepare confirm record and confirm email event") // 500
	}

	// save confirm record and event to send email
	err = uc.repo.CreateConfirm(ctx, &confirmRecord, event.List{confirmEmail})
	if err != nil {
		return errors.Wrap(err, "create confirm record in db") // 500
	}

	return nil
}

func (uc UseCase) prepareConfirmRecordWithEvent(_ context.Context, email, firstName, lastName, lang string) (confirmRecord domain.Confirm, confirmEmail event.Event, err error) {
	const passwordLength = 8
	// prepare variables for email sending
	confirmPassword := utils.RandString(passwordLength)
	confirmRecord, err = domain.NewConfirm(
		domain.ResetEmailPasswordConfirmKind,
		confirmPassword,
		uc.confirmLife,
		map[string]string{ /* vars */
			"email": email,
		},
	)
	if err != nil {
		err = errors.Wrap(err, "create new confirm struct")
		return
	}

	// prepare url
	url := fmt.Sprintf("%s/reset-password?i=%s&v=%s", uc.webBaseURL, utils.UUIDtoB64URL(confirmRecord.ConfirmID), confirmPassword)

	// prepare event email.SendTemplate
	toEmail := email
	toName := fmt.Sprintf("%s %s", firstName, lastName)
	confirmEmail, err = event.EmailSendTemplate(
		resetEmailPasswordTpl,
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

	return confirmRecord, confirmEmail, err
}
