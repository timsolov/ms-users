package postgres

import (
	"database/sql"
	"ms-users/app/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
)

// E helper function to replace specific driver related NotFound error to generic between all db drivers.
// Other errors will be returned without replacing. (gorm.ErrRecordNotFound -> db.ErrNotFound, other error -> other error)
func E(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if IsUniqueViolationErr(err) {
		return domain.ErrNotUnique
	}
	return err
}

func isErr(err error, code string) bool {
	pgerr := &pgconn.PgError{}
	if !errors.As(err, &pgerr) {
		return false
	}
	return pgerr.Code == code
}

// IsUniqueViolationErr returns true if error is unique violation
func IsUniqueViolationErr(err error) bool {
	return isErr(err, pgerrcode.UniqueViolation)
}

// IsForeignKeyViolation returns true if error is foreign key violation
func IsForeignKeyViolation(err error) bool {
	return isErr(err, pgerrcode.ForeignKeyViolation)
}
