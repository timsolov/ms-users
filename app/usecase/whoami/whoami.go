package whoami

import (
	"context"
	"ms-users/app/conf"
	"ms-users/app/domain/entity"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
)

// Repository describes repository contract
type Repository interface {
	// EmailPassIdentByEmail returns email-pass identity by email.
	EmailPassIdentByEmail(ctx context.Context, email string) (ident entity.Ident, err error)
}

// Params describes parameters
type Params struct {
	AccessToken string
}

// UseCase describes usecase
type UseCase struct {
	repo        Repository
	tokenConfig *conf.TOKEN
}

func New(repo Repository, tokenConfig *conf.TOKEN) UseCase {
	return UseCase{
		repo:        repo,
		tokenConfig: tokenConfig,
	}
}

func (uc *UseCase) Do(ctx context.Context, query *Params) (userID uuid.UUID, err error) {
	var jsonToken paseto.JSONToken
	var footer string
	err = paseto.NewV2().Decrypt(query.AccessToken, []byte(uc.tokenConfig.Secret), &jsonToken, &footer)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "decrypt token")
	}

	userID, err = uuid.Parse(jsonToken.Get("user_id"))
	if err != nil {
		err = errors.Wrap(err, "parse uuid from paseto token")
		return
	}

	return
}
