package web

import (
	"ms-users/app/domain/repository"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/infrastructure/logger"
	"ms-users/app/usecase/create_user"
	"ms-users/app/usecase/profile"
)

// Server implements the protobuf interface
type Server struct {
	log logger.Logger
	// queries
	profile profile.ProfileQuery
	// commands
	createUser create_user.CreateUserCommand
}

var _ pb.UserServiceServer = (*Server)(nil)

// New initializes a new Server struct.
func New(log logger.Logger, repo repository.Repository) *Server {
	return &Server{
		// vars
		log: log,
		// queries
		profile: profile.NewProfileQuery(repo),
		// commands
		createUser: create_user.NewCreateUserCommand(repo),
	}
}
