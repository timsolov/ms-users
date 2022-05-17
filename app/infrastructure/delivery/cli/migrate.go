package cli

import (
	"ms-users/app/domain/repository"
	"ms-users/app/infrastructure/logger"
	"os"

	"github.com/spf13/cobra"
)

type migrator interface {
	Migrate(migrateCmd string) (stepsDone int, err error)
}

func NewMigrateCmd(log logger.Logger, r repository.Repository) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate <up|down|1|-1>",
		Short: "Execute migration",
		Long: "Execute migrations stored at binary\n" +
			"Parameters:\n" +
			"up - migrate all steps Up\n" +
			"down - migrate all steps Down\n" +
			"number - amount of steps to migrate (if > 0 - migrate number steps up, if < 0 migrate number steps down)",
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			m, ok := r.(migrator)
			if !ok {
				log.Fatalf("repository desn't support migration")
			}
			log = log.WithFields(logger.Fields{"module": "migrate"})
			n, err := m.Migrate(args[0])
			if err != nil {
				log.Fatalf("migrate: %v", err)
			}
			log.Infof("%d steps done", n)
			os.Exit(0)
		},
	}
}
