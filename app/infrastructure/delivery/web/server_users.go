package web

import (
	"context"

	"ms-users/app/domain/entity"
	"ms-users/app/infrastructure/delivery/web/pb"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/whoami"

	"github.com/google/uuid"
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
		return nil, Unauthorized(ctx)
	}

	user, err := s.profile.Run(ctx, &profile.Params{UserID: userID})
	if err != nil {
		switch errors.Cause(err) {
		case entity.ErrNotFound:
			return out, Unauthorized(ctx)
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

// Whoami returns user_id by access_token.
//
// This end-point considers you have an access_token in Cookie or Authorization header.
// It's possible to use it in authentication middleware for authenticate users.
//
func (s *Server) Whoami(ctx context.Context, _ *pb.WhoamiRequest) (*pb.WhoamiResponse, error) {
	out := &pb.WhoamiResponse{}

	var (
		accessToken string
		userID      uuid.UUID
	)

	// first we check cookie for access_token
	accessToken, err := Cookie(ctx, "access_token")
	if err != nil && err != entity.ErrNotFound {
		return out, Unauthorized(ctx, err)
	}
	if accessToken != "" {
		goto usecase
	}

	// second we look at the bearer token
	accessToken, err = Bearer(ctx)
	if err != nil && err != entity.ErrNotFound {
		return out, Unauthorized(ctx, err)
	}
	if accessToken == "" {
		return out, Unauthorized(ctx, err)
	}

usecase:
	userID, err = s.whoami.Do(ctx, &whoami.Params{AccessToken: accessToken})
	if err != nil {
		return out, Unauthorized(ctx, err)
	}

	out.UserId = userID.String()

	return out, OK(ctx)
}
