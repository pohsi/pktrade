package trade

import (
	"strconv"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/internal/auth"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/pkg/log"
)

type resource struct {
	service Service
	logger  log.Logger
}

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {

	res := resource{service, logger}

	r.Get("/trades", res.query)
	r.Get("/trades/<type>", res.get)
	r.Get("/trades/orders", res.query, res.getOrders)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/trades", res.query, res.makeTrade)
}

// query returns recent 50 trade records for all cards
func (r resource) query(c *routing.Context) error {

	r.logger.Info("Enter trade query")

	records, err := r.service.GetPurchaseOrder(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(records)
}

// query returns recent 50 trade records by card type
func (r resource) get(c *routing.Context) error {

	// r.logger.Infof("Enter trade get")

	// records, err := r.service.Get(c.Request.Context(), c.Param("type"))
	// if err != nil {
	// 	return err
	// }

	// return c.Write(records)
	return nil
}

func (r resource) getOrders(c *routing.Context) error {
	return nil
}

func (r resource) makeTrade(c *routing.Context) error {

	var request CreateOrderRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequestError("")
	}
	user := auth.CurrentUser(c.Request.Context())
	if uid, err := strconv.Atoi(user.GetID()); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequestError("")
	} else {
		request.UserId = uid
	}

	request.UserName = user.GetName()

	order, err := r.service.CreateOrder(c.Request.Context(), request)
	if err != nil {
		return err
	}

	return c.Write(order)
}
