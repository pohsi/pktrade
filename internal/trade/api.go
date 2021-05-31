package trade

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/pkg/log"
)

func RegisterHandler(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

}
