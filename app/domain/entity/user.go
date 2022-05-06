package entity

import (
	"time"

	"github.com/google/uuid"
)

// User describes entity.
type User struct {
	UserID    uuid.UUID
	Email     string
	Password  string // encrypted password
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserList []User
