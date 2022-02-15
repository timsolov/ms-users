package pb

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildBadRequest(errs validation.Errors) error {
	errs = errs.Filter().(validation.Errors)

	if errs == nil {
		return nil
	}

	br := &errdetails.BadRequest{}
	for key, err := range errs {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       key,
			Description: err.Error(),
		}
		br.FieldViolations = append(br.FieldViolations, v)
	}

	st := status.New(codes.InvalidArgument, "invalid request")
	std, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}

	return std.Err()
}
