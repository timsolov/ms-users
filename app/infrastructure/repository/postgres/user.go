package postgres

import (
	"context"
	"time"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateUser creates new profile record
func (d *DB) CreateUser(ctx context.Context, m *entity.User) error {
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
func (d *DB) Profile(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User

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

	return user, entity.ErrNotFound
}

// CreateUserAggregate creates new ident record with user record.
func (d *DB) CreateUserAggregate(ctx context.Context, ua *entity.UserAggregate) error {
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

		return nil
	})

	return E(err)
}
