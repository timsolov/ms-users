package retry_confirm

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	validation = validator.New()
)

// CreateEmailPassIdentityUseCase usecase with necessary method
type CreateEmailPassIdentityUseCase interface {
	// PrepareConfirmEmailRecord creates instance of Confirm struct for confirming email.
	PrepareConfirmEmailRecord(email string) (confirm domain.Confirm, err error)
	// PrepareConfirmEmailEvent creates instance of Event for sending confirmation email.
	PrepareConfirmEmailEvent(email, lang, code string, profile []byte) (confirmEmail event.Event, err error)
}

// Repository describes repository contract
type Repository interface {
	// ReadIdentKind returns ident by ident and kind.
	ReadIdentKind(ctx context.Context, ident string, kind domain.IdentKind) (domain.Ident, error)
	// User returns profile record
	User(ctx context.Context, userID uuid.UUID) (domain.User, error)
	// CreateConfirm creates new confirm record
	CreateConfirm(ctx context.Context, m *domain.Confirm, events []event.Event) error
}

// Params describes parameters
type Params struct {
	Ident string
}

// UseCase describes usecase
type UseCase struct {
	repo      Repository
	emailPass CreateEmailPassIdentityUseCase
}

func New(repo Repository, preparer CreateEmailPassIdentityUseCase) UseCase {
	return UseCase{
		repo:      repo,
		emailPass: preparer,
	}
}

func (uc *UseCase) Do(ctx context.Context, cmd *Params) (err error) {
	// cmd.Ident
	const (
		isEmail = "required,email"
	)
	kind := domain.UnknownIdent
	if validation.VarCtx(ctx, cmd.Ident, isEmail) == nil {
		kind = domain.EmailPassIdent
	}

	switch kind {
	case domain.EmailPassIdent:
		return uc.retryEmailPassConfirm(ctx, cmd.Ident)
	default:
		return domain.ErrUnknownIdent
	}
}

func (uc *UseCase) retryEmailPassConfirm(ctx context.Context, email string) error {
	// find ident with this email
	ident, err := uc.repo.ReadIdentKind(ctx, email, domain.EmailPassIdent)
	if err != nil {
		// 204 - no content when domain.ErrNotFound returned
		if errors.Cause(err) == domain.ErrNotFound {
			err = domain.ErrEmailPassNotFound
		}
		return errors.Wrap(err, "request email-pass ident from db") // 500 - otherwise
	}

	// don't do retry for confirmed identities
	if ident.IdentConfirmed {
		return domain.ErrIdentityConfirmed // 400
	}

	// find user profile
	user, err := uc.repo.User(ctx, ident.UserID)
	if err != nil {
		return errors.Wrap(err, "request user profile from db") // 500
	}

	// prepare new confirm record struct
	confirmRecord, err := uc.emailPass.PrepareConfirmEmailRecord(email)
	if err != nil {
		return errors.Wrap(err, "prepare confirm email record") // 500
	}

	// encode to json object which encoded by base64
	code, err := confirmRecord.ToBase64JSON()
	if err != nil {
		return errors.Wrap(err, "encode confirm struct to base64") // 500
	}

	// prepare event for sending new email for confirmation
	confirmEmailEvent, err := uc.emailPass.PrepareConfirmEmailEvent(email, "en", code, user.Profile)
	if err != nil {
		return errors.Wrap(err, "prepare confirm email event") // 500
	}

	// create confirm record in db and send an event to email service
	err = uc.repo.CreateConfirm(ctx, &confirmRecord, event.List{confirmEmailEvent})
	if err != nil {
		return errors.Wrap(err, "create confirm record in db") // 500
	}

	return nil
}
