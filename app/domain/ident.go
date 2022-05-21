package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdentKind int

const (
	UnknownIdent   IdentKind = 0
	EmailPassIdent IdentKind = 1
)

// Ident describes identity for profile
type Ident struct {
	UserID         uuid.UUID
	Ident          string
	IdentConfirmed bool
	Kind           IdentKind
	Password       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
