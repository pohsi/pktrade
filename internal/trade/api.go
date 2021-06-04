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

	r.Get("/trades/records/<type>", res.queryRecordByType)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Put("/trades/status/<type>", res.queryRecord)
	r.Post("/trades", res.makeTrade)
}

const gerRecordsLimit int = 50

// query returns recent 50 trade records for all cards
func (r resource) queryRecord(c *routing.Context) error {

	queryType, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		return err
	}

	user := auth.CurrentUser(c.Request.Context())
	req := GetRecordsRequest{
		UserName:  user.GetName(),
		Limit:     gerRecordsLimit,
		QueryType: queryType,
	}
	if uid, err := strconv.Atoi(user.GetID()); err != nil {
		return errors.BadRequestError("")
	} else {
		req.UserId = uid
	}

	records, err := r.service.GetRecords(c.Request.Context(), req)
	if err != nil {
		return err
	}

	return c.Write(records)
}

// query returns recent 50 trade records by card type
func (r resource) queryRecordByType(c *routing.Context) error {

	cardType, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		return err
	}

	records, err := r.service.GetRecordsByCardType(c.Request.Context(), GetRecordsByCardTypeRequest{
		CardType: CardType(cardType),
		Limit:    gerRecordsLimit,
	})
	if err != nil {
		return err
	}

	return c.Write(records)
}

func (r resource) makeTrade(c *routing.Context) error {

	var request CreateOrderRequest
	if err := c.Read(&request); err != nil {
		return errors.BadRequestError("")
	}
	user := auth.CurrentUser(c.Request.Context())
	if uid, err := strconv.Atoi(user.GetID()); err != nil {
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
