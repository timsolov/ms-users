package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateIdent creates new ident record
func (d *DB) CreateIdent(ctx context.Context, m *domain.Ident) error {
	if m.UserID == uuid.Nil {
		return fmt.Errorf("user_id: required")
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		now := time.Now()
		m.CreatedAt = now
		m.UpdatedAt = now
	}

	var password sql.NullString
	if m.Password != "" {
		password.String = m.Password
		password.Valid = true
	}

	err := d.execr(ctx, 1,
		`INSERT 
			INTO "idents" (user_id, ident, ident_confirmed, kind, password, created_at, updated_at) 
			VALUES (?,?,?,?,?,?,?)`,
		m.UserID, m.Ident, m.IdentConfirmed, m.Kind, password, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// DeleteIdent deletes user by ident and kind.
func (d *DB) DeleteIdent(ctx context.Context, ident string, kind domain.IdentKind) error {
	query := `DELETE FROM "idents" WHERE ident = ? AND kind = ?`
	return d.execr(ctx, 1, query, ident, kind)
}

// ReadIdent returns ident by unique ident.
func (d *DB) ReadIdent(ctx context.Context, ident string) (domain.Ident, error) {
	var res domain.Ident
	r, err := d.one(ctx,
		`SELECT user_id, ident, ident_confirmed, kind, password, created_at, updated_at
			FROM "idents" WHERE ident = ?`,
		ident)
	if err != nil {
		return res, err
	}

	var password sql.NullString

	err = E(r.Scan(&res.UserID, &res.Ident, &res.IdentConfirmed, &res.Kind, &password, &res.CreatedAt, &res.UpdatedAt))
	res.Password = password.String

	return res, err
}

// IdentKind returns ident by ident and kind.
func (d *DB) ReadIdentKind(ctx context.Context, ident string, kind domain.IdentKind) (domain.Ident, error) {
	var res domain.Ident
	r, err := d.one(ctx,
		`SELECT user_id, ident, ident_confirmed, kind, password, created_at, updated_at
			FROM "idents" WHERE ident = ? AND kind = ?`,
		ident, kind)
	if err != nil {
		return res, err
	}

	var password sql.NullString

	err = E(r.Scan(&res.UserID, &res.Ident, &res.IdentConfirmed, &res.Kind, &password, &res.CreatedAt, &res.UpdatedAt))
	res.Password = password.String

	return res, err
}

// EmailPassIdentByEmail returns email-pass identity by email.
func (d *DB) EmailPassIdentByEmail(ctx context.Context, email string) (ident domain.Ident, err error) {
	return d.ReadIdentKind(ctx, email, domain.EmailPassIdent)
}

// UpdIdent updates ident record.
//
// Errors:
//   - domain.ErrNotFound - when record not found.
func (d *DB) UpdIdent(ctx context.Context, m *domain.Ident, events event.List) error {
	m.UpdatedAt = time.Now()

	err := d.atomic(ctx, func(tx *DB) error {
		err := tx.execr(ctx, 1, `UPDATE idents SET ident_confirmed=?, password=?, updated_at=? WHERE ident=? AND kind=?`, m.IdentConfirmed, m.Password, m.UpdatedAt, m.Ident, m.Kind)
		if err != nil {
			if err == domain.ErrMismatch { // update returns 0 rows affected when record didn't find
				err = domain.ErrNotFound // we've to replace the error with right name
			}
			return errors.Wrap(err, "update ident record")
		}

		err = tx.publishEvents(ctx, events)
		if err != nil {
			return errors.Wrap(err, "publish events")
		}
		return nil
	})

	return err
}
