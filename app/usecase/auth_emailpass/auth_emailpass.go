package auth_emailpass

import (
	"context"
	"ms-users/app/common/password"
	"ms-users/app/conf"
	"ms-users/app/domain"
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
)

// Repository describes repository contract
type Repository interface {
	// EmailPassIdentByEmail returns email-pass identity by email.
	EmailPassIdentByEmail(ctx context.Context, email string) (ident domain.Ident, err error)
}

// Params describes parameters
type Params struct {
	Email    string
	Password string
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

func (uc *UseCase) Do(ctx context.Context, cmd *Params) (accessToken, refreshToken string, err error) {
	ident, err := uc.repo.EmailPassIdentByEmail(ctx, cmd.Email)
	if err != nil {
		err = errors.Wrap(err, "request email-pass identify from db")
		return
	}

	if !password.Verify(ident.Password, cmd.Password) {
		err = domain.ErrUnauthorized
		return
	}

	now := time.Now()
	jsonToken := paseto.JSONToken{
		Issuer:     uc.tokenConfig.Issuer,
		IssuedAt:   now,
		Expiration: now.Add(uc.tokenConfig.AccessLife),
		NotBefore:  now,
	}
	jsonToken.Set("user_id", ident.UserID.String())

	accessToken, err = paseto.NewV2().Encrypt([]byte(uc.tokenConfig.Secret), jsonToken, nil)
	if err != nil {
		err = errors.Wrap(err, "encrypt paseto access token")
		return
	}

	return
}
