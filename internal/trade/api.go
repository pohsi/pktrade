package trade

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/pkg/log"
)

type resource struct {
	service Service
	logger  log.Logger
}

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/trade", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
}

func (r resource) query(c *routing.Context) error {

	return nil
}
