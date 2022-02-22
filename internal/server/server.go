package server

import (
	"github.com/timsolov/ms-users/internal/usecase"
)

type Server interface {
	UserServiceServer
}

// server implements the protobuf interface
type server struct {
	uc usecase.UseCase
}

// New initializes a new Backend struct.
func New(uc usecase.UseCase) Server {
	return &server{
		uc: uc,
	}
}
