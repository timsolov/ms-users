package cli

import (
	"ms-users/app/infrastructure/logger"
	"ms-users/app/usecase/create_emailpass_identity"
	"os"

	"github.com/spf13/cobra"
)

func NewCreateEmailPassIdentityCmd(log logger.Logger, uc create_emailpass_identity.UseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "create-email-pass-identity <email> <password> <firstname> <lastname>",
		Short: "Create user and email-pass identity for him",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			userID, err := uc.Run(cmd.Context(), &create_emailpass_identity.Params{
				Email:          args[0],
				EmailConfirmed: true,
				Password:       args[1],
				FirstName:      args[2],
				LastName:       args[3],
			})
			if err != nil {
				log.Fatalf("create-email-pass-identity error: %v", err)
			}
			log.Infof("new profile created with ID: %s", userID)
			os.Exit(0)
		},
	}
}
