package accesslog

import (
	"context"
	"net/http"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/access"
	"github.com/google/uuid"
	"github.com/pohsi/pktrade/pkg/log"
)

const (
	xRequestIdField     = "X-Request-ID"
	xCorrelationIdField = "X-Correlation-ID"
)

func getRequestId(req *http.Request) string {
	return req.Header.Get(xRequestIdField)
}

func getCorrelationId(req *http.Request) string {
	return req.Header.Get(xCorrelationIdField)
}

func withRequest(ctx context.Context, req *http.Request) context.Context {
	id := getRequestId(req)
	if id == "" {
		id = uuid.New().String()
	}

	ctx = context.WithValue(ctx, log.RequestIdKey, id)
	if cid := getCorrelationId(req); id != "" {
		ctx = context.WithValue(ctx, log.CorrelatationIdKey, cid)
	}

	return ctx
}

func NewHandler(logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {

		start := time.Now()

		rw := &access.LogResponseWriter{ResponseWriter: c.Response, Status: http.StatusOK}
		c.Response = rw

		ctx := c.Request.Context()
		ctx = withRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)
		err := c.Next()

		logger.With(ctx, "duration", time.Now().Sub(start).Milliseconds(), "status",
			rw.Status).Infof("%s %s %s %d %d", c.Request.Method, c.Request.URL.Path, c.Request.Proto,
			rw.Status, rw.BytesWritten)

		return err
	}
}
