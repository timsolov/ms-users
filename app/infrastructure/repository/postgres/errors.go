package postgres

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
)

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
