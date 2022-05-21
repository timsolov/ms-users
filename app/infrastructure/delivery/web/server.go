package web

import (
	"ms-users/app/conf"
	"ms-users/app/domain/repository"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/infrastructure/logger"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/whoami"
)

// Server implements the protobuf interface
type Server struct {
	log logger.Logger
	// queries
	profile profile.UseCase
	whoami  whoami.UseCase
	// commands
	createEmailPassIdentity create_emailpass_identity.UseCase
	authEmailPass           auth_emailpass.UseCase
}

var _ pb.UserServiceServer = (*Server)(nil)

// New initializes a new Server struct.
func New(log logger.Logger, repo repository.Repository, config *conf.Config) *Server {
	return &Server{
		// vars
		log: log,
		// queries
		profile: profile.New(repo),
		whoami:  whoami.New(repo, &config.TOKEN),
		// commands
		createEmailPassIdentity: create_emailpass_identity.New(repo),
		authEmailPass:           auth_emailpass.New(repo, &config.TOKEN),
	}
}
