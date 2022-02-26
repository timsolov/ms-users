package web

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthRequired = status.Error(codes.Unauthenticated, "authorization required")
)

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}
