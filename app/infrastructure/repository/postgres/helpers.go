package postgres

import (
	"context"
	"database/sql"

	"ms-users/app/domain/entity"
)

func (d *DB) execr(ctx context.Context, rows int64, query string, args ...interface{}) error {
	r, err := d.db.ExecContext(ctx, d.db.Rebind(query), args...)

	if err != nil {
		return err
	}

	if rows > 0 {
		i, err := r.RowsAffected()
		if err != nil {
			return err
		}
		if i != rows {
			return entity.ErrMismatch
		}
	}

	return nil
}

func E(err error) error {
	if err == sql.ErrNoRows {
		return entity.ErrNotFound
	}
	return err
}
