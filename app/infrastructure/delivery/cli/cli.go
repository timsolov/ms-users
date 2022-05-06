package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func Run(cmds ...*cobra.Command) error {
	rootCmd := &cobra.Command{Use: "app"}

	// switch off usage message on run without args
	rootCmd.Run = func(cmd *cobra.Command, args []string) {}

	// add exit on help
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		defaultHelp(c, s)
		os.Exit(-1)
	})

	// add commands
	rootCmd.AddCommand(cmds...)

	return rootCmd.Execute()
}
