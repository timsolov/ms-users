package postgres

import (
	"context"
	"ms-users/app/domain"
)

// Publish saves to PgQ an event.
// This function can't be tested properly.
func (d *DB) Publish(ctx context.Context, subject string, payload []byte) error {
	row, err := d.one(ctx, "SELECT pgq.insert_event('outbox', ?, ?)", subject, payload)
	if err != nil {
		return err
	}

	// the insert_event function should returns Event ID
	// see: https://pgq.github.io/extension/pgq/files/external-sql.html#pgq.insert_event(3)
	var r int64
	err = row.Scan(&r)
	if err != nil {
		return E(err)
	}

	if r == 0 {
		return domain.ErrMismatch
	}

	return nil
}
