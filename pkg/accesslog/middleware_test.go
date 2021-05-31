package accesslog

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

const exampleComUrl = "http://example.com"

func buildRequest(requestId, correlationId string) (*http.Request, error) {
	req, err := http.NewRequest("GET", exampleComUrl, bytes.NewBufferString(""))

	if requestId != "" {
		req.Header.Set(xRequestIdField, requestId)
	}

	if correlationId != "" {
		req.Header.Set(xCorrelationIdField, correlationId)
	}

	return req, err

}

func Test_getRequestId(t *testing.T) {
	req, _ := buildRequest("", "")
	assert.Empty(t, getRequestId(req))
	req.Header.Set(xRequestIdField, "test")
	assert.Equal(t, "test", getRequestId(req))
}

func Test_getCorrelationId(t *testing.T) {
	req, _ := buildRequest("", "")
	assert.Empty(t, getCorrelationId(req))
	req.Header.Set(xCorrelationIdField, "test")
	assert.Equal(t, "test", getCorrelationId(req))
}

func Test_withRequest(t *testing.T) {
	req, _ := buildRequest("aaa", "777")
	ctx := withRequest(context.Background(), req)
	assert.Equal(t, "aaa", ctx.Value(log.RequestIdKey).(string))
	assert.Equal(t, "777", ctx.Value(log.CorrelatationIdKey).(string))

	req, _ = buildRequest("", "777")
	ctx = withRequest(context.Background(), req)

	// If RequestIdKey is not exist, generate new one
	assert.NotEmpty(t, ctx.Value(log.RequestIdKey).(string))
	assert.Equal(t, "777", ctx.Value(log.CorrelatationIdKey).(string))
}

func TestNewHandler(t *testing.T) {

	res := httptest.NewRecorder()
	req, err := http.NewRequest("Get", "http://127.0.0.1/users", nil)
	assert.Equal(t, true, req != nil && err == nil)
	ctx := routing.NewContext(res, req)

	logger, entries := log.NewForTest()
	handler := NewHandler(logger)
	err = handler(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, entries.Len())
	assert.Equal(t, "Get /users HTTP/1.1 200 0", entries.All()[0].Message)

}
