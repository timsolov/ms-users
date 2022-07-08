package grpc_server

import (
	"context"
	"net/http"
	"strconv"

	"ms-users/app/common/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Errorf builds a gRPC status error
func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

// 400
func BadRequest(ctx context.Context, err error) error {
	return Custom(ctx, codes.InvalidArgument, http.StatusBadRequest, err)
}

// 204
func NoContent(ctx context.Context, err ...error) error {
	var e error
	if len(err) == 1 {
		e = err[0]
	}
	return Custom(ctx, codes.NotFound, http.StatusNoContent, e)
}

// 404
func NotFound(ctx context.Context, err ...error) error {
	var e error
	if len(err) == 1 {
		e = err[0]
	}
	return Custom(ctx, codes.NotFound, http.StatusNotFound, e)
}

// 403
func Forbidden(ctx context.Context, err ...error) error {
	var e error
	if len(err) == 1 {
		e = err[0]
	}
	return Custom(ctx, codes.PermissionDenied, http.StatusForbidden, e)
}

// 401
func Unauthorized(ctx context.Context, err ...error) error {
	var e error
	if len(err) == 1 {
		e = err[0]
	}
	return Custom(ctx, codes.PermissionDenied, http.StatusUnauthorized, e)
}

// 500
func Internal(ctx context.Context, log logger.Logger, format string, args ...interface{}) error {
	log.Errorf(format, args...)
	return Custom(ctx, codes.Internal, http.StatusInternalServerError, nil)
}

// Custom code
func Custom(ctx context.Context, code codes.Code, statusCode int, err error) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(statusCode)))
	if err != nil {
		return status.Error(code, err.Error())
	}
	return status.Error(code, "")
}

// 200
func OK(ctx context.Context) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "200"))
	return nil
}
