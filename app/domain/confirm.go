package domain

import (
	"encoding/base64"
	"encoding/json"
	"ms-users/app/common/password"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ConfirmKind int

const (
	UnknownConfirmKind = 0
	EmailConfirmKind   = 1
)

// Confirm describes confirm body
type Confirm struct {
	ConfirmID string            `json:"c"`
	Password  string            `json:"p"`
	Kind      ConfirmKind       `json:"-"`
	Vars      map[string]string `json:"-"`
	CreatedAt time.Time         `json:"-"`
	ValidTill time.Time         `json:"-"`
}

func NewConfirm(kind ConfirmKind, life time.Duration, vars map[string]string) (confirm Confirm, err error) {
	confirmPass, err := password.Encrypt(uuid.NewString())
	if err != nil {
		err = errors.Wrap(err, "encrypt confirm password")
		return
	}
	confirm.ConfirmID = uuid.NewString()
	confirm.Password = confirmPass
	confirm.CreatedAt = time.Now()
	confirm.ValidTill = confirm.CreatedAt.Add(life)
	confirm.Kind = kind
	confirm.Vars = vars
	return
}

func (c *Confirm) ToBase64() (encoded string, err error) {
	confirmJSON, err := json.Marshal(c)
	if err != nil {
		err = errors.Wrap(err, "marshal confirm struct")
		return
	}
	encoded = base64.URLEncoding.EncodeToString(confirmJSON)
	return
}

func (c *Confirm) FromBase64(encoded string) (err error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		err = errors.Wrap(err, "decode from base64")
		return
	}

	err = json.Unmarshal(decoded, c)
	if err != nil {
		err = errors.Wrap(err, "unmarshal confirm struct")
		return
	}

	return
}
