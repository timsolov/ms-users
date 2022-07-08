package grpc_server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	validation = validator.New()
)

func Validate(ctx context.Context, v interface{}) error {
	err := validation.StructCtx(ctx, v)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)

	br := &errdetails.BadRequest{}
	for _, elem := range errs {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       elem.Field(),
			Description: elem.Error(),
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
