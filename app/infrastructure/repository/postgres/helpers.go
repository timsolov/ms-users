package postgres

import (
	"context"
	"database/sql"

	"ms-users/app/domain/entity"
)

func (d *DB) execr(ctx context.Context, rows int64, query string, args ...interface{}) error {
	var (
		r   sql.Result
		err error
	)

	if d.tx != nil {
		r, err = d.tx.ExecContext(ctx, d.db.Rebind(query), args...)
	} else {
		r, err = d.tx.ExecContext(ctx, d.db.Rebind(query), args...)
	}
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
