package web

import (
	"context"

	"github.com/pkg/errors"
	"github.com/timsolov/ms-users/app/domain/entity"
	"github.com/timsolov/ms-users/app/pb"
	"github.com/timsolov/ms-users/app/usecase/create_user"
	"github.com/timsolov/ms-users/app/usecase/profile"
)

// stub: s *Server pb.UserServiceServer

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
func (s *Server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// nolint:gocritic
	// userID := grpc_gateway.UserID(ctx)
	// if userID == "" {
	// 	return &pb.CreateUserResponse{}, ErrAuthRequired
	// }

	if stErr := in.Validate(); stErr != nil {
		return &pb.CreateUserResponse{}, stErr
	}

	createUser := create_user.CreateUser{
		Email:     in.GetEmail(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
	}

	userID, err := s.createUser.Do(ctx, &createUser)
	if err != nil {
		return &pb.CreateUserResponse{}, Internal(ctx, "usecase NewUser: %s", err)
	}

	return &pb.CreateUserResponse{
		UserId: userID.String(),
	}, nil
}

// Profile detail info.
//
// Profile returns user detail info.
func (s *Server) Profile(ctx context.Context, _ *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	userID, err := XUserId(ctx)
	if err != nil {
		return nil, Forbidden(ctx)
	}

	user, err := s.profile.Do(ctx, &profile.Profile{UserID: userID})
	if err != nil {
		switch errors.Cause(err) {
		case entity.ErrNotFound:
			return &pb.ProfileResponse{}, BadRequest(ctx, err)
		default:
			return &pb.ProfileResponse{}, Internal(ctx, "")
		}
	}

	return &pb.ProfileResponse{
		UserId:    user.UserID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
