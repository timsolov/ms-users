package grpc_server

import (
	"context"

	"ms-users/app/common/logger"
	"ms-users/app/delivery/grpc_server/pb"
	"ms-users/app/domain"
	"ms-users/app/usecase/auth_emailpass"
	"ms-users/app/usecase/confirm"
	"ms-users/app/usecase/create_emailpass_identity"
	"ms-users/app/usecase/profile"
	"ms-users/app/usecase/reset_password_confirm"
	"ms-users/app/usecase/reset_password_init"
	"ms-users/app/usecase/retry_confirm"
	"ms-users/app/usecase/update_profile"
	"ms-users/app/usecase/whoami"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// stub: s *Server pb.UserServiceServer

// Creates new profile with email-password identity.
//
// Access: Public
//
// For creating new profile you have to provide:
// - email (required);
// - password (required);
// - first_name (optional);
// - last_name (optional);
func (s *Server) CreateEmailPassIdentity(ctx context.Context, in *pb.CreateEmailPassIdentityRequest) (*pb.CreateEmailPassIdentityResponse, error) {
	out := &pb.CreateEmailPassIdentityResponse{}

	if stErr := Validate(ctx, in); stErr != nil {
		return out, stErr
	}

	profileJSON, err := in.Profile.MarshalJSON()
	if err != nil {
		return out, BadRequest(ctx, errors.Wrap(err, "profile"))
	}

	createUser := create_emailpass_identity.Params{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Profile:  profileJSON,
	}

	userID, err := s.commands.CreateEmailPassIdentity.Do(ctx, &createUser)
	if err != nil {
		_, isStatus := status.FromError(err)
		if isStatus {
			return out, err
		}

		switch errors.Cause(err) {
		case domain.ErrIdentityDuplicated:
			err = BadRequest(ctx, err)
		default:
			err = Internal(ctx, s.log, "CreateEmailPassIdentity usecase: %s", err)
		}

		return out, err
	}

	out.UserId = userID.String()

	return out, OK(ctx)
}

// Profile detail info.
//
// Access: X-User-Id
//
// Profile returns profile and identities of user by user_id.
func (s *Server) Profile(ctx context.Context, _ *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	out := &pb.ProfileResponse{}

	userID, err := XUserId(ctx)
	if err != nil {
		return nil, Unauthorized(ctx)
	}

	user, err := s.queries.Profile.Do(ctx, &profile.Params{UserID: userID})
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrNotFound:
			return out, Unauthorized(ctx)
		default:
			return out, Internal(ctx, s.log, "Profile usecase: %s", err)
		}
	}

	out.UserId = user.UserID.String()
	out.Profile, err = structpb.NewStruct(nil)
	if err != nil {
		return out, Internal(ctx, s.log, "init profile field")
	}
	err = protojson.Unmarshal(user.Profile, out.Profile)
	if err != nil {
		return out, Internal(ctx, s.log, "unmarshal profile field from json to proto")
	}
	if len(user.Idents) > 0 {
		out.Idents = make([]*pb.Identity, 0, len(user.Idents))
		for i := 0; i < len(user.Idents); i++ {
			ident := &pb.Identity{
				Ident: user.Idents[i].Ident,
				Kind:  pb.Identity_Kind(user.Idents[i].Kind),
			}

			out.Idents = append(out.Idents, ident)
		}
	}

	return out, OK(ctx)
}

// Confirm universal confirm link
//
// Access: Public
//
// It's possible to confirm different type of operations.
func (s *Server) Confirm(ctx context.Context, in *pb.ConfirmRequest) (out *pb.ConfirmResponse, err error) {
	out = &pb.ConfirmResponse{}

	err = s.commands.Confirm.Do(ctx, &confirm.Params{
		Encoded: in.GetEncoded(),
	})
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrNotFound, domain.ErrMismatch, domain.ErrExpired:
			return out, BadRequest(ctx, err)
		default:
			return out, Internal(ctx, s.log, "Confirm usecase: %s", err)
		}
	}

	return out, OK(ctx)
}

// RetryConfirm resends confirmation code.
//
// Access: Public
//
// This end-point is utilized when confirmation code is expired and
// user wants to reissue new confirmation code.
//
// For email-pass identity should be provided email and if identity with related
// email exists confirmation will be sent to that email.
func (s *Server) RetryConfirm(ctx context.Context, in *pb.RetryConfirmRequest) (*pb.RetryConfirmResponse, error) {
	out := &pb.RetryConfirmResponse{}

	ident := in.GetIdent()

	// we will do validation inside usecase
	err := s.commands.RetryConfirm.Do(ctx, &retry_confirm.Params{
		Ident: ident,
	})

	switch errors.Cause(err) {
	case nil:
		// pass
	case domain.ErrEmailPassNotFound, domain.ErrIdentityConfirmed: // 200
		// pass
	case domain.ErrUnknownIdent: // 400
		return out, BadRequest(ctx, err)
	default:
		return out, Internal(ctx, s.log, "RetryConfirm usecase: %s", err)
	}

	return out, OK(ctx)
}

// Authenticate users by email-password.
//
// Access: Public
func (s *Server) AuthEmailPass(ctx context.Context, in *pb.AuthEmailPassRequest) (*pb.AuthEmailPassResponse, error) {
	out := &pb.AuthEmailPassResponse{}

	accessToken, _, err := s.commands.AuthEmailPass.Do(ctx, &auth_emailpass.Params{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrNotConfirmed, domain.ErrUnauthorized:
			err = Unauthorized(ctx, err)
		default:
			err = Internal(ctx, s.log, "AuthEmailPass usecase: %s", err)
		}
		return out, err
	}

	out.AccessToken = accessToken

	return out, OK(ctx)
}

// Whoami returns user_id by access_token.
//
// Access: Bearer or Cookie
//
// This end-point considers you have an access_token in Cookie or Authorization header.
// It's possible to use it in authentication middleware for authenticate users.
func (s *Server) Whoami(ctx context.Context, _ *pb.WhoamiRequest) (*pb.WhoamiResponse, error) {
	out := &pb.WhoamiResponse{}

	var (
		accessToken string
		userID      uuid.UUID
	)

	// first we check cookie for access_token
	accessToken, err := Cookie(ctx, "access_token")
	if err != nil && err != domain.ErrNotFound {
		return out, Unauthorized(ctx, err)
	}
	if accessToken != "" {
		goto usecase
	}

	// second we look at the bearer token
	accessToken, err = Bearer(ctx)
	if err != nil && err != domain.ErrNotFound {
		return out, Unauthorized(ctx, err)
	}
	if accessToken == "" {
		return out, Unauthorized(ctx, err)
	}

usecase:
	userID, err = s.queries.Whoami.Do(ctx, &whoami.Params{AccessToken: accessToken})
	if err != nil {
		return out, Unauthorized(ctx, err)
	}

	out.UserId = userID.String()

	return out, OK(ctx)
}

// ResetPasswordInit begins reset-password process for identity.
//
// Access: Public
//
// For email-pass identity should be provided email and if identity with related
// email exists confirmation code for recovery process will be sent to that email.
// In email will be stored link with comfirm_id (i) and verifycation code (p) inside
// query parameters of the link. It should leads the user to the web page where the
// user will see input for `new password`.
//
// This end-point will always return 200 OK for failed and success requests. This is
// necessary to prevent database brute-forcing.
func (s *Server) ResetPasswordInit(ctx context.Context, in *pb.ResetPasswordInitRequest) (*pb.ResetPasswordInitResponse, error) {
	out := &pb.ResetPasswordInitResponse{}

	ident := in.GetIdent()

	err := s.commands.ResetPasswordInit.Do(ctx, &reset_password_init.Params{Ident: ident})
	if err != nil {
		level := logger.ErrorLevel
		switch errors.Cause(err) {
		case domain.ErrUnknownIdent:
			return out, BadRequest(ctx, domain.ErrUnknownIdent)
		case domain.ErrEmailPassNotFound:
			level = logger.WarnLevel
		}
		s.log.Logf(level, "ResetPasswordInit usecase: %s", err)
	}

	return out, OK(ctx)
}

// ResetPasswordConfirm confirm identity recovery process and set new password.
//
// Access: Public
//
// It's necessary to identify does the user who started recovery process is owner of
// the identity. So this end-point waits for verification id, code and new password.
func (s *Server) ResetPasswordConfirm(ctx context.Context, in *pb.ResetPasswordConfirmRequest) (*pb.ResetPasswordConfirmResponse, error) {
	out := &pb.ResetPasswordConfirmResponse{}

	if stErr := Validate(ctx, in); stErr != nil {
		return out, stErr
	}

	err := s.commands.ResetPasswordConfirm.Do(ctx, &reset_password_confirm.Params{
		ConfirmIDB64: in.GetConfirmId(),
		Verification: in.GetVerifycation(),
		NewPassword:  in.GetPassword(),
	})
	if err != nil {
		switch errors.Cause(err) {
		case domain.ErrNotFound, domain.ErrExpired, domain.ErrMismatch:
			return out, BadRequest(ctx, err)
		}
		return out, Internal(ctx, s.log, "ResetPasswordConfirm usecase: %s", err)
	}

	return out, OK(ctx)
}

// UpdateProfile updates profile traits.
//
// Access: X-User-Id
//
// Updates one or multiple profile traits in database.
func (s *Server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	out := &pb.UpdateProfileResponse{}

	userID, err := XUserId(ctx)
	if err != nil {
		return nil, Unauthorized(ctx)
	}

	if stErr := Validate(ctx, in); stErr != nil {
		return out, stErr
	}

	profileJSON, err := in.Profile.MarshalJSON()
	if err != nil {
		return out, BadRequest(ctx, errors.Wrap(err, "profile"))
	}

	err = s.commands.UpdateProfile.Do(ctx, &update_profile.Params{
		UserID:  userID,
		Profile: profileJSON,
	})
	if err != nil {
		_, isStatus := status.FromError(err)
		if isStatus {
			return out, err
		}

		err = Internal(ctx, s.log, "UpdateProfile usecase: %s", err)

		return out, err
	}

	return out, OK(ctx)
}
