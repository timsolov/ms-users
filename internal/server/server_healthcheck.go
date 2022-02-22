package server

import "net/http"

// TODO: write healthcheck
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
