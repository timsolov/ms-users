package web

import (
	"ms-users/app/domain/repository"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/infrastructure/logger"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
)

// Server implements the protobuf interface
type Server struct {
	log logger.Logger
	// queries
	profile profile.UseCase
	// commands
	createUserPassIdentity create_emailpass_identity.UseCase
}

var _ pb.UserServiceServer = (*Server)(nil)

// New initializes a new Server struct.
func New(log logger.Logger, repo repository.Repository) *Server {
	return &Server{
		// vars
		log: log,
		// queries
		profile: profile.New(repo),
		// commands
		createUserPassIdentity: create_emailpass_identity.New(repo),
	}
}
