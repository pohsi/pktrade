package errors

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

func handlerResponseOK(c *routing.Context) error {
	return c.Write("ok")
}

func handlerResponseError(c *routing.Context) error {
	return fmt.Errorf("error")
}

func handlerResponseHTTPError(c *routing.Context) error {
	return notFoundError("not found")
}

func handlerPanic(c *routing.Context) error {
	panic("panic")
}

func buildContext(handlers ...routing.Handler) (*routing.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://127.0.0.1/users", nil)
	return routing.NewContext(res, req, handlers...), res
}

func TestNewHandler(t *testing.T) {
	t.Run("normal processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := NewHandler(logger)
		ctx, res := buildContext(handler, handlerResponseOK)
		assert.Nil(t, ctx.Next())
		assert.Zero(t, entries.Len())
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("error processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := NewHandler(logger)
		ctx, res := buildContext(handler, handlerResponseError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 1, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("http error processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := NewHandler(logger)
		ctx, res := buildContext(handler, handlerResponseHTTPError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("panic processing", func(t *testing.T) {
		logger, entries := log.NewForTest()
		handler := NewHandler(logger)
		ctx, res := buildContext(handler, handlerPanic)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 2, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func Test_buildErrorResponce(t *testing.T) {
	err := notFoundError("")
	assert.Equal(t, err, buildErrorResponse(err))

	err = buildErrorResponse(routing.NewHTTPError(http.StatusNotFound))
	assert.Equal(t, http.StatusNotFound, err.Status)

	err = buildErrorResponse(validation.Errors{})
	assert.Equal(t, http.StatusBadRequest, err.Status)

	err = buildErrorResponse(routing.NewHTTPError(http.StatusForbidden))
	assert.Equal(t, http.StatusForbidden, err.Status)

	err = buildErrorResponse(sql.ErrNoRows)
	assert.Equal(t, http.StatusNotFound, err.Status)

	err = buildErrorResponse(fmt.Errorf(("test")))
	assert.Equal(t, http.StatusInternalServerError, err.Status)
}
