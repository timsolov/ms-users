package reset_password_confirm

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/common/password"
	"ms-users/app/common/utils"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Repository describes repository contract
type Repository interface {
	// IdentKind returns ident by ident and kind.
	ReadIdentKind(ctx context.Context, ident string, kind domain.IdentKind) (domain.Ident, error)
	// ReadConfirm returns confirm record by confirm_id.
	ReadConfirm(ctx context.Context, confirmID uuid.UUID) (confirm domain.Confirm, err error)
	// UpdIdent updates ident record.
	UpdIdent(ctx context.Context, m *domain.Ident, events event.List) error
	// DelConfirm deletes confirm record.
	DelConfirm(ctx context.Context, confirmID uuid.UUID) (err error)
}

// Params describes parameters
type Params struct {
	ConfirmIDB64 string
	Verification string
	NewPassword  string
}

// UseCase describes usecase
type UseCase struct {
	repo Repository
}

// New creates new instance of the UseCase
func New(repo Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

// Do the main method for the UseCase
func (uc UseCase) Do(ctx context.Context, cmd *Params) (err error) {
	confirmID, err := utils.B64URLtoUUID(cmd.ConfirmIDB64)
	if err != nil {
		return errors.Wrap(err, "convert confirm id from base64url to uuid") // 400
	}

	// looking for confirm record by confirm_id
	confirm, err := uc.repo.ReadConfirm(ctx, confirmID)
	if err != nil {
		// domain.ErrNotFound - 400
		return errors.Wrap(err, "request confirm record from db") // other - 500
	}

	// verify request expiration
	if confirm.ValidTill.Before(time.Now()) {
		return domain.ErrExpired // 400
	}

	// check does public password is same to stored encrypted
	if !password.Verify(confirm.Password, cmd.Verification) { // not same
		return domain.ErrMismatch // 400
	}

	var (
		ident     string
		identKind domain.IdentKind
	)
	switch confirm.Kind {
	case domain.ResetEmailPasswordConfirmKind:
		var ok bool
		ident, ok = confirm.Vars["email"]
		if !ok {
			return domain.ErrBadFormat // 500
		}
		identKind = domain.EmailPassIdent
	default:
		return domain.ErrUnknownConfirmKind // 500
	}

	identRecord, err := uc.repo.ReadIdentKind(ctx, ident, identKind)
	if err != nil {
		// domain.ErrNotFound - 500
		if err == domain.ErrNotFound {
			err = domain.ErrIdentNotFound // 500
		}
		return errors.Wrap(err, "request ident record from db") // 500
	}

	identRecord.Password, err = password.Encrypt(cmd.NewPassword)
	if err != nil {
		return errors.Wrap(err, "encrypt password") // 500
	}

	// save new password
	err = uc.repo.UpdIdent(ctx, &identRecord, event.None)
	if err != nil {
		return errors.Wrap(err, "update ident record in db") // 500
	}

	_ = uc.repo.DelConfirm(ctx, confirmID)

	return nil
}
