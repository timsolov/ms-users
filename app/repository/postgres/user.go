package postgres

import (
	"context"
	"fmt"
	"ms-users/app/common/event"
	"ms-users/app/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateUser creates new user record
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

var (
	columnsUpdUser       = []string{"view", "profile"}
	columnsUpdRestricted = []string{"created_at", "updated_at"}
)

// UpdUser updates user record.
func (d *DB) UpdUser(ctx context.Context, m *domain.User, columns ...string) error {
	if m.UserID == uuid.Nil {
		return fmt.Errorf("user_id: required")
	}

	if len(columns) == 0 {
		columns = columnsUpdUser
	}

	if mcontains(columns, columnsUpdRestricted, atLeastOne) {
		return fmt.Errorf("restricted columns for update: %v", columnsUpdRestricted)
	}

	const two = 2
	conditions := make([]string, 0, len(columns)+two) // +2 because we'll add "created_at" and "user_id"
	args := make([]any, 0, len(columns)+two)

	conditions = append(conditions, "updated_at=?")
	args = append(args, time.Now())

	for _, column := range columns {
		switch column {
		case "view":
			conditions = append(conditions, "view=?")
			args = append(args, m.View)
		case "profile":
			conditions = append(conditions, "profile=?")
			args = append(args, m.Profile)
		}
	}

	args = append(args, m.UserID)

	query := fmt.Sprintf(`UPDATE "users" SET %s WHERE user_id = ?`,
		strings.Join(conditions, ","),
	)

	err := d.execr(ctx, 1,
		query,
		args...)
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

// User returns profile record
func (d *DB) User(ctx context.Context, userID uuid.UUID) (domain.User, error) {
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

		if confirm != nil {
			err = tx.CreateConfirm(ctx, confirm, event.None)
			if err != nil {
				return errors.Wrap(err, "create confirm")
			}
		}

		err = tx.publishEvents(ctx, events)
		if err != nil {
			return errors.Wrap(err, "publish events")
		}

		return nil
	})

	return E(err)
}
