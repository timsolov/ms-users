package web

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*emptypb.Empty, error) {
	// userID := grpc_gateway.UserID(ctx)
	// if userID == "" {
	// 	return &emptypb.Empty{}, ErrAuthRequired
	// }

	errs := validation.Errors{
		"email": validation.Validate(in.Email, validation.Required, is.Email),
	}
	if err := errs.Filter(); err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
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
		return &emptypb.Empty{}, Errorf(codes.Internal, "usecase NewUser: %s", err)
	}

	// _ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "204"))

	return &emptypb.Empty{}, nil
}

// List users.
//
// Returns the list of users records.
// Maximum records per request is 100.
// Pagination available by using offset, limit.
func (s *server) ListUsers(_ context.Context, query *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	panic("not impelemented")
}

// Update user info.
//
// Update user info fully or partial.
func (s *server) UpdateUser(_ context.Context, in *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	panic("not impelemented")
}

// UserDetail detail info.
//
// UserDetail returns user detail info.
func (s *server) UserDetail(context.Context, *pb.UserRequest) (*pb.User, error) {
	panic("not impelemented")
}
