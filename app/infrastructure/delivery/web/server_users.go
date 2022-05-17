package web

import (
	"context"

	"ms-users/app/domain/entity"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/usecase/create_emailpass_identity"
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
func (s *Server) CreateEmailPassIdentity(ctx context.Context, in *pb.CreateEmailPassIdentityRequest) (*pb.CreateEmailPassIdentityResponse, error) {
	if stErr := Validate(ctx, in); stErr != nil {
		return &pb.CreateEmailPassIdentityResponse{}, stErr
	}

	createUser := create_emailpass_identity.Params{
		Email:     in.GetEmail(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Password:  in.GetPassword(),
	}

	userID, err := s.createUserPassIdentity.Run(ctx, &createUser)
	if err != nil {
		switch errors.Cause(err) {
		case entity.ErrNotUnique:
			err = BadRequest(ctx, ErrIdentityDuplicated)
		default:
			err = Internal(ctx, s.log, "usecase createUserPassIdentity: %s", err)
		}
		return &pb.CreateEmailPassIdentityResponse{}, err
	}

	return &pb.CreateEmailPassIdentityResponse{
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

	user, err := s.profile.Run(ctx, &profile.Params{UserID: userID})
	if err != nil {
		switch errors.Cause(err) {
		case entity.ErrNotFound:
			return &pb.ProfileResponse{}, Forbidden(ctx)
		default:
			return &pb.ProfileResponse{}, Internal(ctx, s.log, "usecase: %s", err)
		}
	}

	var profileView entity.V1Profile
	err = user.UnmarshalProfile(&profileView)
	if err != nil {
		return &pb.ProfileResponse{}, Internal(ctx, s.log, "unmarshaling profile: %s", err)
	}

	return &pb.ProfileResponse{
		UserId:    user.UserID.String(),
		Email:     profileView.Email,
		FirstName: profileView.FirstName,
		LastName:  profileView.LastName,
	}, OK(ctx)
}

// Confirm universal confirm link
//
// It's possible to confirm different type of operations.
//
func (s *Server) Confirm(_ context.Context, _ *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	panic("not implemented") // TODO: Implement
}
