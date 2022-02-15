package postgres

import (
	"embed"
	"strconv"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*
var fs embed.FS

// Migrate does migration of database.
//   - migrateCmd - up/down, number/-number steps.
func (d *DB) Migrate(migrateCmd string) (stepsDone int, err error) {
	var steps int
	const max = 999

	switch migrateCmd {
	case "up":
		steps = max
	case "down":
		steps = -max
	default:
		var err error
		steps, err = strconv.Atoi(migrateCmd)
		if err != nil {
			return 0, errors.Wrap(err, "failed to convert migrate argument to digit")
		}
	}

	// Migrator
	if steps != 0 {
		migrations := &migrate.EmbedFileSystemMigrationSource{
			FileSystem: fs,
		}

		var direction migrate.MigrationDirection

		if steps > 0 {
			direction = migrate.Up
		} else if steps < 0 {
			direction = migrate.Down
			steps = -steps
		}

		db, err := d.SqlDB()
		if err != nil {
			return 0, errors.Wrap(err, "failed to get SqlDB")
		}

		n, err := migrate.ExecMax(db, "postgres", migrations, direction, steps)
		if err != nil {
			return 0, errors.Wrap(err, "failed to execute migrations")
		}

		if direction == migrate.Down {
			n = -n
		}

		return n, nil
	}

	return 0, nil
}
