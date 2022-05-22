package postgres

import (
	"context"
	"ms-users/app/domain"
	"time"
)

// CreateConfirm creates new confirm record
func (d *DB) CreateConfirm(ctx context.Context, m *domain.Confirm) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	err := d.execr(ctx, 1,
		`INSERT 
			INTO "confirms" (confirm_id, password, kind, vars, created_at, valid_till)
			VALUES (?,?,?,?,?,?)`,
		m.ConfirmID, m.Password, m.Kind, m.Vars, m.CreatedAt, m.ValidTill)
	if err != nil {
		return err
	}

	return nil
}
