package web

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthRequired = status.Error(codes.Unauthenticated, "authorization required")
)

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

func BadRequest(ctx context.Context, err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}

func Internal(ctx context.Context, format string, args ...interface{}) error {
	//TODO: write log to console
	return status.Error(codes.Internal, "internal")
}

func Custom(ctx context.Context, statusCode int, err error) error {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(statusCode)))
	return status.Error(codes.Unknown, err.Error())
}
