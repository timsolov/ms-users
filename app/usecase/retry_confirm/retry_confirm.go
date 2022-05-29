package retry_confirm

import (
	"context"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

var (
	validation = validator.New()
)

// CreateEmailPassIdentityUseCase usecase with necessary method
type CreateEmailPassIdentityUseCase interface {
	PrepareConfirmRecordAndConfirmEmail(firstName, lastName, email, lang string) (confirmRecord domain.Confirm, confirmEmail event.Event, err error)
}

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

	// find user profile
	profile, err := uc.repo.Profile(ctx, ident.UserID)
	if err != nil {
		return errors.Wrap(err, "request user profile from db") // 500
	}

	firstName := gjson.GetBytes(profile.Profile, "first_name").String()
	lastName := gjson.GetBytes(profile.Profile, "last_name").String()
	if firstName == "" || lastName == "" {
		return fmt.Errorf("first_name or last_name is empty") // 500
	}

	confirmRecord, confirmEmail, err := uc.emailPass.PrepareConfirmRecordAndConfirmEmail(firstName, lastName, email, "en")
	if err != nil {
		return errors.Wrap(err, "prepare confirm record and confirm email event") // 500
	}

	err = uc.repo.CreateConfirm(ctx, &confirmRecord, event.List{confirmEmail})
	if err != nil {
		return errors.Wrap(err, "create confirm record in db") // 500
	}

	return nil
}
