package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
)

// CreateUser creates new user record
func (d *DB) CreateUser(ctx context.Context, m *entity.User) error {
	if m.UserID == uuid.Nil {
		m.UserID = uuid.New()
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		now := time.Now()
		m.CreatedAt = now
		m.UpdatedAt = now
	}

	err := d.execr(ctx, 1,
		`INSERT 
			INTO "users" (user_id, email, first_name, last_name, created_at, updated_at) 
			VALUES (?,?,?,?,?,?)`,
		m.UserID, m.Email, m.FirstName, m.LastName, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Profile returns user record
func (d *DB) Profile(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User

	query := "SELECT user_id, email, first_name, last_name, created_at, updated_at FROM users WHERE user_id = ?"

	rows, err := d.db.QueryContext(ctx, d.db.Rebind(query), userID)
	if err != nil {
		return user, E(err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
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
