package cli

import (
	"context"
	"ms-users/app/infrastructure/logger"
	"ms-users/app/usecase/create_user"
	"os"

	"github.com/spf13/cobra"
)

func NewCreateUserCmd(log logger.Logger, createUser create_user.CreateUserCommand) *cobra.Command {
	return &cobra.Command{
		Use:   "create-user <email> <password> <firstname> <lastname>",
		Short: "create-user",
		Long:  "create-user",
		Args:  cobra.ExactArgs(4),
		Run: func(_ *cobra.Command, args []string) {
			userID, err := createUser.Do(context.Background(), &create_user.CreateUser{
				Email:     args[0],
				FirstName: args[2],
				LastName:  args[3],
				Password:  args[1],
			})
			if err != nil {
				log.Fatalf("create-user error: %v", err)
			}
			log.Infof("new user created with ID: %s", userID)
			os.Exit(0)
		},
	}
}
