package cli

import (
	"ms-users/app/common/logger"
	"ms-users/app/usecase/create_emailpass_identity"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"google.golang.org/grpc/status"
)

func NewCreateEmailPassIdentityCmd(log logger.Logger, uc *create_emailpass_identity.UseCase) *cobra.Command {
	const (
		minArgs  = 2
		emailIdx = 0
		passIdx  = 1
	)
	return &cobra.Command{
		Use:     "create-email-pass-identity <email> <password> <profile field=value> ...<other profile fields>",
		Short:   "Create user and email-pass identity for him",
		Example: "./service create-email-pass-identity user@example.org pa55w0rd first_name=John last_name=Doe",
		Args:    cobra.MinimumNArgs(minArgs),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				profile []byte
				err     error
			)

			// fill out profile
			for i := 2; i < len(args); i++ {
				const twoParts = 2

				parts := strings.SplitN(args[i], "=", twoParts)
				if len(parts) != 2 {
					continue
				}

				profile, err = sjson.SetBytes(profile, parts[0], parts[1])
				if err != nil {
					log.Fatalf("set profile field: %s=%s", parts[0], parts[1])
				}
			}

			userID, err := uc.Do(cmd.Context(), &create_emailpass_identity.Params{
				Email:          args[emailIdx],
				EmailConfirmed: true,
				Password:       args[passIdx],
				Profile:        profile,
			})
			if err != nil {
				s, isStatus := status.FromError(err)
				if isStatus {
					log.Fatalf("%s", s.Details())
				}
				log.Fatalf("create-email-pass-identity error: %s", err)
			}
			log.Infof("new profile created with ID: %s", userID)
			os.Exit(0)
		},
	}
}
