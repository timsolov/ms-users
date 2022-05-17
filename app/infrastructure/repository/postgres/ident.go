package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ms-users/app/domain/entity"

	"github.com/google/uuid"
)

// CreateIdent creates new ident record
func (d *DB) CreateIdent(ctx context.Context, m *entity.Ident) error {
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
