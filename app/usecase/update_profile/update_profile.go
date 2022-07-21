package update_profile

import (
	"context"
	"ms-users/app/common/jsonschema"
	"ms-users/app/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type JSONSchemaValidator interface {
	// ValidateProfile validates profile object by jsonschema and build gRPC status error with bad request.
	ValidateProfile(ctx context.Context, profile []byte) error
}

// Repository describes repository contract
type Repository interface {
	// UpdUser updates user record.
	UpdUser(ctx context.Context, m *domain.User, columns ...string) error
}

// Params describes parameters
type Params struct {
	UserID  uuid.UUID
	Profile []byte // json
}

// UseCase describes usecase
type UseCase struct {
	repo                Repository
	jsonSchema          jsonschema.Schema
	jsonSchemaValidator JSONSchemaValidator
}

func New(
	repo Repository,
	jsonSchema jsonschema.Schema,
	jsValidator JSONSchemaValidator,
) UseCase {
	uc := UseCase{
		repo:                repo,
		jsonSchema:          jsonSchema,
		jsonSchemaValidator: jsValidator,
	}

	return uc
}

func (uc *UseCase) Do(ctx context.Context, cmd *Params) (err error) {
	err = uc.jsonSchemaValidator.ValidateProfile(ctx, cmd.Profile)
	if err != nil {
		return err // 400 (status error) or 500 (not status error)
	}

	err = uc.repo.UpdUser(ctx, &domain.User{UserID: cmd.UserID, Profile: cmd.Profile}, "profile")
	if err != nil {
		return errors.Wrap(err, "update user profile in db")
	}

	return nil
}
