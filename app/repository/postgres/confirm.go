package postgres

import (
	"context"
	"ms-users/app/common/event"
	"ms-users/app/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateConfirm creates new confirm record
func (d *DB) CreateConfirm(ctx context.Context, m *domain.Confirm, events []event.Event) error {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	return d.atomic(ctx, func(tx *DB) error {
		err := d.execr(ctx, 1,
			`INSERT 
				INTO "confirms" (confirm_id, password, kind, vars, created_at, valid_till)
				VALUES (?,?,?,?,?,?)`,
			m.ConfirmID, m.EncryptedPassword, m.Kind, m.Vars, m.CreatedAt, m.ValidTill)
		if err != nil {
			return errors.Wrap(err, "create confirm record")
		}

		err = tx.publishEvents(ctx, events)
		if err != nil {
			return errors.Wrap(err, "publish events")
		}

		return nil
	})
}

// ReadConfirm returns confirm record by confirm_id.
func (d *DB) ReadConfirm(ctx context.Context, confirmID uuid.UUID) (confirm domain.Confirm, err error) {
	r, err := d.one(ctx, `SELECT confirm_id, password, kind, vars, created_at, valid_till FROM confirms WHERE confirm_id = ?`, confirmID)
	if err != nil {
		return
	}

	var vars StringInterfaceMap

	err = r.Scan(&confirm.ConfirmID, &confirm.Password, &confirm.Kind, &vars, &confirm.CreatedAt, &confirm.ValidTill)
	if err != nil {
		err = E(err)
		return
	}

	if len(vars) > 0 {
		confirm.Vars = make(map[string]string)
		for k, v := range vars {
			if s, ok := v.(string); ok {
				confirm.Vars[k] = s
			} else {
				err = domain.ErrBadFormat
				return
			}
		}
	}

	return confirm, nil
}

// DelConfirm deletes confirm record.
func (d *DB) DelConfirm(ctx context.Context, confirmID uuid.UUID) (err error) {
	return E(d.execr(ctx, 1, `DELETE FROM confirms WHERE confirm_id = ?`, confirmID))
}
