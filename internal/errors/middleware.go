package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pohsi/pktrade/pkg/log"
)

func NewHandler(logger log.Logger) routing.Handler {
	return func(c *routing.Context) (err error) {
		defer func() {
			l := logger.With(c.Request.Context())
			if e := recover(); e != nil {
				var ok bool
				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}
				l.Errorf("recorvered from panic (%v): %s", err, debug.Stack())
			}

			if err != nil {
				res := buildErrorResponse(err)
				if res.StatusCode() == http.StatusInternalServerError {
					l.Errorf("encountered internal server error: %v", err)
				}
				c.Response.WriteHeader(res.StatusCode())
				if err = c.Write(res); err != nil {
					l.Errorf("failed writing error response: v", err)
				}
				c.Abort()
				err = nil
			}
		}()
		return c.Next()
	}
}

func buildErrorResponse(err error) errorResponse {
	switch err.(type) {
	case errorResponse:
		return err.(errorResponse)
	case validation.Errors:
		return invalidInput(err.(validation.Errors))
	case routing.HTTPError:
		switch err.(routing.HTTPError).StatusCode() {
		case http.StatusNotFound:
			return notFoundError("")
		default:
			return errorResponse{
				Status:  err.(routing.HTTPError).StatusCode(),
				Message: err.Error(),
			}
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return notFoundError("")
	}

	return internalServerError("")
}
