package web

import (
	"context"

	"ms-users/app/domain/entity"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"

	"github.com/pkg/errors"
)

// stub: s *Server pb.UserServiceServer

// Creates new user.
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

	userID, err := s.createEmailPassIdentity.Run(ctx, &createUser)
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
	out := &pb.ProfileResponse{}

	userID, err := XUserId(ctx)
	if err != nil {
		return nil, Forbidden(ctx)
	}

	user, err := s.profile.Run(ctx, &profile.Params{UserID: userID})
	if err != nil {
		switch errors.Cause(err) {
		case entity.ErrNotFound:
			return out, Forbidden(ctx)
		default:
			return out, Internal(ctx, s.log, "usecase: %s", err)
		}
	}

	var profileView entity.V1Profile
	err = user.UnmarshalProfile(&profileView)
	if err != nil {
		return out, Internal(ctx, s.log, "unmarshaling profile: %s", err)
	}

	out.UserId = user.UserID.String()
	out.Email = profileView.Email
	out.FirstName = profileView.FirstName
	out.LastName = profileView.LastName

	return out, OK(ctx)
}

// Confirm universal confirm link
//
// It's possible to confirm different type of operations.
//
func (s *Server) Confirm(_ context.Context, _ *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	panic("not implemented") // TODO: Implement
}

// Authenticate users by email-pasword.
//
// Access: any
//
func (s *Server) AuthEmailPass(ctx context.Context, in *pb.AuthEmailPassRequest) (*pb.AuthEmailPassResponse, error) {
	out := &pb.AuthEmailPassResponse{}

	accessToken, _, err := s.authEmailPass.Do(ctx, &auth_emailpass.Params{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		return out, Internal(ctx, s.log, "usecase: %s", err)
	}

	out.AccessToken = accessToken

	return out, OK(ctx)
}
