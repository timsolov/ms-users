package web

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthRequired = status.Error(http.StatusUnauthorized, "authorization required")
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

func Internal(ctx context.Context, format string, args ...interface{}) error {
	//TODO: write log to console
	return Custom(ctx, codes.Internal, http.StatusInternalServerError, nil)
}

func Custom(ctx context.Context, code codes.Code, statusCode int, err error) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(statusCode)))
	return status.Error(code, err.Error())
}
