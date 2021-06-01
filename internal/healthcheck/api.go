package healthcheck

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

func RegisterHandlers(r *routing.Router, response string) {
	r.To("GET,HEAD", "/healthcheck", healthyMessage(response))
}

func healthyMessage(msg string) routing.Handler {
	return func(c *routing.Context) error {
		return c.Write("OK: " + msg)
	}
}
