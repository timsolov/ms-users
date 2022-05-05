package web

import (
	"context"

	"ms-users/app/domain/entity"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/usecase/create_user"
	"ms-users/app/usecase/profile"

	"github.com/pkg/errors"
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
		return &pb.CreateUserResponse{}, Internal(ctx, s.log, "usecase createUser: %s", err)
	}

	return &pb.CreateUserResponse{
		UserId: userID.String(),
	}, OK(ctx)
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
			return &pb.ProfileResponse{}, Forbidden(ctx)
		default:
			return &pb.ProfileResponse{}, Internal(ctx, s.log, "usecase profile: %s", err)
		}
	}

	return &pb.ProfileResponse{
		UserId:    user.UserID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, OK(ctx)
}
