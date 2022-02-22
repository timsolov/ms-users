package gateway

import (
	"context"

	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/metadata"
)

// UserID returns user id from metadata.
// UserID should be UUID formated user's identificator passed to a Header `X-User-Id`
// for RESTAPI request or to metadata `x-user-id` for gRPC request.
func UserID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if uID, ok := md["x-user-id"]; ok &&
			len(uID) == 1 &&
			is.UUID.Validate(uID[0]) == nil {
			return uID[0]
		}
	}
	return ""
}
