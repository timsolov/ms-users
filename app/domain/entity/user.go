package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User describes entity.
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

// V1Profile describes view "v1"
type V1Profile struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserAggregate describes user aggregate model.
type UserAggregate struct {
	User
	Idents []Ident // one user can have many identities
}
