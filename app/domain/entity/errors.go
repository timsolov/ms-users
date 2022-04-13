package entity

import "errors"

// Universal errors
var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrNotUnique = errors.New("not unique")
	ErrMismatch  = errors.New("mismatch")
)
