package event

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/dimiro1/health"
)

func HealthChecker(url string) health.CheckerFunc {
	return func(ctx context.Context) health.Health {
		var res health.Health
		res.Down()

		start := time.Now()

		response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{})) // nolint: gosec
		if err != nil {
			res.AddInfo("error", err.Error())
			return res
		}
		defer response.Body.Close()

		duration := time.Since(start)
		res.AddInfo("duration", duration.String())

		if response.StatusCode == http.StatusOK {
			res.Up()
		}

		return res
	}
}
