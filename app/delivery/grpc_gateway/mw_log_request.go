package grpc_gateway

import (
	"net/http"
	"time"

	"ms-users/app/common/logger"
	"ms-users/app/common/rw"
)

func LogRequest(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const healthcheck = "/healthcheck"

			started := time.Now()
			rwCustom := &rw.ResponseWriter{ResponseWriter: w}

			next.ServeHTTP(rwCustom, r)

			var level logger.Level
			switch {
			case rwCustom.Code >= 500: //nolint
				level = logger.ErrorLevel
			case rwCustom.Code >= 400: //nolint
				level = logger.WarnLevel
			case r.RequestURI == healthcheck: // remove healthcheck from logs
				return
			default:
				level = logger.InfoLevel
			}

			log.Logf(
				level,
				"%d %s %s (%v)",
				rwCustom.Code,
				r.Method,
				r.RequestURI,
				time.Since(started),
			)
		})
	}
}
