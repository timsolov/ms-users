package web

import (
	"context"
	"ms-users/app/domain/entity"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// XUserId checks if set X-User-Id header and it has valid UUID value.
func XUserId(ctx context.Context) (userID uuid.UUID, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if uID, ok := md["x-user-id"]; ok && len(uID) == 1 {
			return uuid.Parse(uID[0])
		}
	}
	err = entity.ErrNotFound
	return
}

// Cookie returns value of cookie with name `name`.
// Returns entity.ErrNotFound when cookie is absent.
func Cookie(ctx context.Context, name string) (value string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = entity.ErrNotFound
		return
	}
	const cookie = "cookie"
	if allCookie, ok := md[cookie]; ok && len(allCookie) > 0 {
		header := http.Header{}
		for i := 0; i < len(allCookie); i++ {
			header.Add("Cookie", allCookie[i])
		}
		request := http.Request{Header: header}

		var c *http.Cookie

		c, err = request.Cookie(name)
		if err != nil {
			if err == http.ErrNoCookie {
				err = entity.ErrNotFound
			}
			return
		}

		value = c.Value
		return
	}

	err = entity.ErrNotFound
	return
}

// Bearer extracts bearer token from Authorization header.
// Returns entity.ErrNotFound when token is absent.
func Bearer(ctx context.Context) (value string, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authorization, ok := md["authorization"]; ok && len(authorization) > 0 {
			const partsAmount = 2
			parts := strings.SplitN(authorization[0], " ", partsAmount)
			if len(parts) != partsAmount {
				err = entity.ErrNotFound
				return
			}
			value = strings.TrimSuffix(parts[1], ";")
			return
		}
	}
	err = entity.ErrNotFound
	return
}
