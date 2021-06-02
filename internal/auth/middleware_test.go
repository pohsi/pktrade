package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestCurrentUser(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, CurrentUser(ctx))
	ctx = WithUser(ctx, "100", "test")
	identity := CurrentUser(ctx)
	if assert.NotNil(t, identity) {
		assert.Equal(t, "100", identity.GetID())
		assert.Equal(t, "test", identity.GetName())
	}
}

func TestHandler(t *testing.T) {
	assert.NotNil(t, NewHandler("test"))
}

func Test_handleToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx, _ := test.MockRoutingContext(req)
	assert.Nil(t, CurrentUser(ctx.Request.Context()))

	err := handleToken(ctx, &jwt.Token{
		Claims: jwt.MapClaims{
			"id":   "100",
			"name": "test",
		},
	})
	assert.Nil(t, err)
	identity := CurrentUser(ctx.Request.Context())
	if assert.NotNil(t, identity) {
		assert.Equal(t, "100", identity.GetID())
		assert.Equal(t, "test", identity.GetName())
	}
}

func TestMocks(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx, _ := test.MockRoutingContext(req)
	assert.NotNil(t, MockAuthHandler(ctx))
	req.Header = MockAuthHeader()
	ctx, _ = test.MockRoutingContext(req)
	assert.Nil(t, MockAuthHandler(ctx))
	assert.NotNil(t, CurrentUser(ctx.Request.Context()))
}

func MockAuthHandler(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "TEST" {
		return errors.UnauthorizedError("")
	}
	ctx := WithUser(c.Request.Context(), "100", "Tester")
	c.Request = c.Request.WithContext(ctx)
	return nil
}

func MockAuthHeader() http.Header {
	header := http.Header{}
	header.Add("Authorization", "TEST")
	return header
}
