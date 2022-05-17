package web

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"ms-users/app/infrastructure/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrIdentityDuplicated = errors.New("identity duplicated")
)

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

func BadRequest(ctx context.Context, err error) error {
	return Custom(ctx, codes.InvalidArgument, http.StatusBadRequest, err)
}

func Forbidden(ctx context.Context) error {
	return Custom(ctx, codes.PermissionDenied, http.StatusForbidden, nil)
}

func Internal(ctx context.Context, log logger.Logger, format string, args ...interface{}) error {
	log.Errorf(format, args...)
	return Custom(ctx, codes.Internal, http.StatusInternalServerError, nil)
}

func Custom(ctx context.Context, code codes.Code, statusCode int, err error) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(statusCode)))
	if err != nil {
		return status.Error(code, err.Error())
	}
	return status.Error(code, "")
}

func OK(ctx context.Context) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "200"))
	return nil
}
