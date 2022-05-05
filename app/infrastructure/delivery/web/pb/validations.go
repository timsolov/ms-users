package pb

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate validates CreateUserRequest and returns status.Error.
func (r *CreateUserRequest) Validate() error {
	errs := validation.Errors{
		"email":      validation.Validate(r.GetEmail(), validation.Required, is.Email),
		"password":   validation.Validate(r.GetPassword(), validation.Required),
		"first_name": validation.Validate(r.GetFirstName(), validation.Required),
		"last_name":  validation.Validate(r.GetLastName(), validation.Required),
	}

	return buildBadRequest(errs)
}
