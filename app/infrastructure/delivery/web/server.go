package web

import (
	"github.com/timsolov/ms-users/app/domain/repository"
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase/create_user"
)

type Server interface {
	pb.UserServiceServer
}

// server implements the protobuf interface
type server struct {
	createUser create_user.CreateUserCommand
}

// New initializes a new Backend struct.
func New(repo repository.Repository) Server {
	return &server{
		createUser: create_user.NewCreateUserCommand(repo),
	}
}
