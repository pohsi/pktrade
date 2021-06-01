package auth

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/pkg/log"
)

func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.Logger) {
	rg.Post("/login", login(service, logger))
}

func login(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errors.BadRequestError("")
		}

		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}
