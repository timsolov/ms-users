package web

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/pb"
)

// stub: s *server UserServiceServer

// Creates new user.
//
// For creating new user you have to provide:
// - email;
// - password;
// ```json
// {
//   "name": "value"
// }
// ```
func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// nolint:gocritic
	// userID := grpc_gateway.UserID(ctx)
	// if userID == "" {
	// 	return &pb.CreateUserResponse{}, ErrAuthRequired
	// }

	if stErr := in.Validate(); stErr != nil {
		return &pb.CreateUserResponse{}, stErr
	}

	user := entity.User{
		UserID:    uuid.New(),
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.uc.NewUser(ctx, &user)
	if err != nil {
		return &pb.CreateUserResponse{}, Internal(ctx, "usecase NewUser: %s", err)
	}

	return &pb.CreateUserResponse{}, nil
}

// List users.
//
// Returns the list of users records.
// Maximum records per request is 100.
// Pagination available by using offset, limit.
func (s *server) ListUsers(_ context.Context, query *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	panic("not implemented")
}

// Update user info.
//
// Update user info fully or partial.
func (s *server) UpdateUser(_ context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	panic("not implemented")
}

// UserDetail detail info.
//
// UserDetail returns user detail info.
func (s *server) UserDetail(context.Context, *pb.UserDetailRequest) (*pb.UserDetailResponse, error) {
	panic("not implemented")
}
