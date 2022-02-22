package db

import (
	"fmt"
)

// migrator interface describes Migrate func.
type migrator interface {
	// Migrate does migration of database.
	//   - migrateCmd - up/down, number/-number steps.
	Migrate(migrateCmd string) (stepsDone int, err error)
}

// Migrate detects Migrator method from db.DB interface.
//   - migrateCmd - command for migration one of:
//      - up/down - do all steps up or down;
//      - number/-number - do amount of steps up when positive or down if negative number.
func Migrate(d DB, migrateCmd string) (stepsDone int, err error) {
	migrator, ok := d.(migrator)
	if !ok {
		return 0, fmt.Errorf("migrator not supported")
	}
	return migrator.Migrate(migrateCmd)
}
