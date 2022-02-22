package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User describes entity.
type User struct {
	UserID    uuid.UUID
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserList []User

type UserRepository interface {
	// NewUser creates new user record
	NewUser(ctx context.Context, m *User) error
	// User returns user record by id
	User(ctx context.Context, userID uuid.UUID, columns ...string) (User, error)
	// UpdUser changes user record's properties.
	UpdUser(ctx context.Context, m *User, columns ...string) error
	// DelUser deletes user record record.
	DelUser(ctx context.Context, userID uuid.UUID) error
}

type UserUseCase interface {
	// NewUser creates new user record
	NewUser(ctx context.Context, user *User) error
}
