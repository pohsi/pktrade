package test

import (
	"net/http"
	"net/http/httptest"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/pkg/accesslog"
	"github.com/pohsi/pktrade/pkg/log"
)

func MockRoutingContext(req *http.Request) (*routing.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	ctx := routing.NewContext(res, req)
	ctx.SetDataWriter(&content.JSONDataWriter{})
	return ctx, res
}

func MockRouter(logger log.Logger) *routing.Router {
	router := routing.New()
	router.Use(
		accesslog.NewHandler(logger),
		errors.NewHandler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)
	return router
}
