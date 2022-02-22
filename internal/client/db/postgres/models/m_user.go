package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/timsolov/ms-users/internal/entity"
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

// ToEntity converts to entity
func (m *User) ToEntity() entity.User {
	return entity.User{
		UserID:    m.UserID,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity converts from entity
func (m *User) FromEntity(e *entity.User) {
	m.UserID = e.UserID
	m.Email = e.Email
	m.FirstName = e.FirstName
	m.LastName = e.LastName
	m.CreatedAt = e.CreatedAt
	m.UpdatedAt = e.UpdatedAt
}

// UserList list
type UserList []User

// ToEntity converts to entity
func (l UserList) ToEntity() entity.UserList {
	if len(l) == 0 {
		return nil
	}
	r := make(entity.UserList, len(l))
	for i := 0; i < len(l); i++ {
		r[i] = l[i].ToEntity()
	}
	return r
}
