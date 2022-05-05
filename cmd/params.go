package main

import (
	"flag"
	"os"

	"ms-users/app/domain/repository"
	"ms-users/app/infrastructure/logger"
)

type migrator interface {
	Migrate(migrateCmd string) (stepsDone int, err error)
}

func ParseParams(log logger.Logger, r repository.Repository) {
	var (
		migrateCmd string
	)

	flag.StringVar(&migrateCmd, "migrate", "",
		"up - migrate all steps Up\n"+
			"down - migrate all steps Down\n"+
			"number - amount of steps to migrate (if > 0 - migrate number steps up, if < 0 migrate number steps down)")
	flag.Parse()

	if len(migrateCmd) > 0 {
		m, ok := r.(migrator)
		if !ok {
			log.Fatalf("repository desn't support migration")
		}
		log = log.WithFields(logger.Fields{"module": "migrate"})
		n, err := m.Migrate(migrateCmd)
		if err != nil {
			log.Fatalf("migrate: %v", err)
		}
		log.Infof("%d steps done", n)
		os.Exit(0)
	}
}
