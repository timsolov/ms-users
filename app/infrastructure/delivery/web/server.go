package web

import (
	"github.com/timsolov/ms-users/app/domain/repository"
	"github.com/timsolov/ms-users/app/infrastructure/logger"
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase/create_user"
	"github.com/timsolov/ms-users/app/usecase/profile"
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
