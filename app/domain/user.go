package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User describes domain.
type User struct {
	UserID    uuid.UUID
	View      string
	Profile   []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) UnmarshalProfile(v interface{}) error {
	return json.Unmarshal(u.Profile, v)
}

func (u *User) MarshalProfile(v interface{}) error {
	var err error
	u.Profile, err = json.Marshal(v)
	return err
}

// UserAggregate describes user aggregate model.
type UserAggregate struct {
	User
	Idents []Ident // one user can have many identities
}
