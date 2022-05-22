package grpc_gateway

import (
	"net/http"
	"time"

	"ms-users/app/common/logger"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LogRequest(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const healthcheck = "/healthcheck"

			started := time.Now()
			rw := &responseWriter{w, 0}

			next.ServeHTTP(rw, r)

			var level logger.Level
			switch {
			case rw.code >= 500: //nolint
				level = logger.ErrorLevel
			case rw.code >= 400: //nolint
				level = logger.WarnLevel
			case r.RequestURI == healthcheck: // remove healthcheck from logs
				return
			default:
				level = logger.InfoLevel
			}

			log.Logf(
				level,
				"%d %s %s (%v)",
				rw.code,
				r.Method,
				r.RequestURI,
				time.Since(started),
			)
		})
	}
}
