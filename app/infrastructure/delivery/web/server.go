package web

import (
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase"
)

type Server interface {
	pb.UserServiceServer
}

// server implements the protobuf interface
type server struct {
	uc usecase.Usecase
}

// New initializes a new Backend struct.
func New(uc usecase.Usecase) Server {
	return &server{
		uc: uc,
	}
}
