package web

import (
	"ms-users/app/common/logger"
	"ms-users/app/delivery/web/pb"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/confirm"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/reset_password_confirm"
	"ms-users/app/usecase/reset_password_init"
	"ms-users/app/usecase/retry_confirm"
	"ms-users/app/usecase/whoami"
)

// Queries describes usecases
type Queries struct {
	Profile profile.UseCase
	Whoami  whoami.UseCase
}

// Commands describes usecases
type Commands struct {
	CreateEmailPassIdentity create_emailpass_identity.UseCase
	AuthEmailPass           auth_emailpass.UseCase
	Confirm                 confirm.UseCase
	RetryConfirm            retry_confirm.UseCase
	ResetPasswordInit       reset_password_init.UseCase
	ResetPasswordConfirm    reset_password_confirm.UseCase
}

// Server implements the protobuf interface
type Server struct {
	log logger.Logger

	queries  *Queries
	commands *Commands
}

var _ pb.UserServiceServer = (*Server)(nil)

// New initializes a new Server struct.
func New(log logger.Logger, queries *Queries, commands *Commands) *Server {
	return &Server{
		// vars
		log: log,
		// queries
		queries: queries,
		// commands
		commands: commands,
	}
}
