package web

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// XUserId checks if set X-User-Id header and it has valid UUID value.
func XUserId(ctx context.Context) (uuid.UUID, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if uID, ok := md["x-user-id"]; ok && len(uID) == 1 {
			return uuid.Parse(uID[0])
		}
	}
	return uuid.Nil, nil
}
