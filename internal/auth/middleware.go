package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/pohsi/pktrade/internal/entity"
)

func NewHandler(verificationKey string) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
}

func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(int),
		token.Claims.(jwt.MapClaims)["name"].(string),
	)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

func WithUser(ctx context.Context, id int, name string) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, Name: name})
}

func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}
	return nil
}
