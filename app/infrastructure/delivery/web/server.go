package web

import (
	"github.com/timsolov/ms-users/app/domain/repository"
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase/create_user"
)

// Server implements the protobuf interface
type Server struct {
	createUser create_user.CreateUserCommand
}

var _ pb.UserServiceServer = (*Server)(nil)

// New initializes a new Server struct.
func New(repo repository.Repository) *Server {
	return &Server{
		createUser: create_user.NewCreateUserCommand(repo),
	}
}
