package cli

import (
	"ms-users/app/common/logger"
	"ms-users/app/usecase/create_emailpass_identity"
	"os"

	"github.com/spf13/cobra"
)

func NewCreateEmailPassIdentityCmd(log logger.Logger, uc create_emailpass_identity.UseCase) *cobra.Command {
	const (
		exactArgs    = 4
		emailIdx     = 0
		passIdx      = 1
		firstNameIdx = 2
		lastNameIdx  = 3
	)
	return &cobra.Command{
		Use:   "create-email-pass-identity <email> <password> <firstname> <lastname>",
		Short: "Create user and email-pass identity for him",
		Args:  cobra.ExactArgs(exactArgs),
		Run: func(cmd *cobra.Command, args []string) {
			userID, err := uc.Do(cmd.Context(), &create_emailpass_identity.Params{
				Email:          args[emailIdx],
				EmailConfirmed: true,
				Password:       args[passIdx],
				FirstName:      args[firstNameIdx],
				LastName:       args[lastNameIdx],
			})
			if err != nil {
				log.Fatalf("create-email-pass-identity error: %v", err)
			}
			log.Infof("new profile created with ID: %s", userID)
			os.Exit(0)
		},
	}
}
