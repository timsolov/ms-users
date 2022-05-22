package postgres

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateUser creates new profile record
func (d *DB) CreateUser(ctx context.Context, m *domain.User) error {
	if m.UserID == uuid.Nil {
		m.UserID = uuid.New()
	}

	now := time.Now()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}

	err := d.execr(ctx, 1,
		`INSERT 
			INTO "users" (user_id, view, profile, created_at, updated_at)
			VALUES (?,?,?,?,?)`,
		m.UserID, m.View, m.Profile, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes user by id.
func (d *DB) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM users WHERE user_id = ?"
	return d.execr(ctx, 1, query, userID)
}

// Profile returns profile record
func (d *DB) Profile(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	var user domain.User

	query := "SELECT user_id, view, profile, created_at, updated_at FROM users WHERE user_id = ?"

	rows, err := d.db.QueryContext(ctx, d.db.Rebind(query), userID)
	if err != nil {
		return user, E(err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.View, &user.Profile, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return user, E(err)
		}
		return user, nil
	}
	err = rows.Err()
	if err != nil {
		return user, E(err)
	}

	return user, domain.ErrNotFound
}

// CreateUserAggregate creates new ident record with user record.
func (d *DB) CreateUserAggregate(ctx context.Context, ua *domain.UserAggregate, confirm *domain.Confirm, events []event.Event) error {
	err := d.atomic(ctx, func(tx *DB) error {
		err := tx.CreateUser(ctx, &ua.User)
		if err != nil {
			return errors.Wrap(err, "create user")
		}

		for i := 0; i < len(ua.Idents); i++ {
			err = tx.CreateIdent(ctx, &ua.Idents[i])
			if err != nil {
				return errors.Wrapf(err, "create ident[%d]", i)
			}
		}

		for i := 0; i < len(events); i++ {
			err = tx.Publish(ctx, events[i].Subject, events[i].Payload)
			if err != nil {
				return errors.Wrapf(err, "publish events[%d]", i)
			}
		}

		if confirm != nil {
			err = tx.CreateConfirm(ctx, confirm)
			if err != nil {
				return errors.Wrap(err, "create confirm")
			}
		}

		return nil
	})

	return E(err)
}
