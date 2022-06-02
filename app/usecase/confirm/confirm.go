package confirm

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
	// ReadConfirm returns confirm record by confirm_id.
	ReadConfirm(ctx context.Context, confirmID uuid.UUID) (confirm domain.Confirm, err error)
	// EmailPassIdentByEmail returns email-pass identity by email.
	EmailPassIdentByEmail(ctx context.Context, email string) (ident domain.Ident, err error)
	// UpdIdent updates ident record.
	//
	// Errors:
	//   - domain.ErrNotFound - when record not found.
	//
	UpdIdent(ctx context.Context, m *domain.Ident, events event.List) error
}

// Params describes parameters
type Params struct {
	Encoded string
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

func (uc UseCase) Do(ctx context.Context, cmd *Params) (err error) {
	var confirmPublic domain.Confirm

	// decode encoded parameter
	err = confirmPublic.FromBase64(cmd.Encoded)
	if err != nil {
		return err // 500
	}

	confirm, err := uc.repo.ReadConfirm(ctx, confirmPublic.ConfirmID)
	if err != nil {
		return errors.Wrap(err, "read confirm record from db") // 400 - not found. 500 - other
	}

	// check expiration
	if confirm.ValidTill.Before(time.Now()) {
		return domain.ErrExpired // 400
	}

	// check does public password is same to stored encrypted
	if !password.Verify(confirm.Password, confirmPublic.Password) { // not same
		return domain.ErrMismatch // 400
	}

	switch confirm.Kind {
	case domain.EmailConfirmKind:
		email, ok := confirm.Vars["email"]
		if !ok {
			return domain.ErrNotFound // 400
		}
		err = uc.confirmEmailIdentity(ctx, email)
		if err != nil {
			return errors.Wrap(err, "confirm email identity") // 400 - not found, 500 - other
		}
	default:
		return fmt.Errorf("unknown confirm kind")
	}

	return nil
}

func (uc UseCase) confirmEmailIdentity(ctx context.Context, email string) error {
	ident, err := uc.repo.EmailPassIdentByEmail(ctx, email)
	if err != nil {
		return errors.Wrap(err, "request email-pass ident from db")
	}

	ident.IdentConfirmed = true

	err = uc.repo.UpdIdent(ctx, &ident, event.None)
	if err != nil {
		return errors.Wrap(err, "update ident record in db")
	}

	return nil
}
