package auth

import (
	"context"
	"testing"

	"github.com/pohsi/pktrade/internal/entity"
	"github.com/pohsi/pktrade/internal/errors"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

func Test_service_Authenticate(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService("test", 100, logger)
	_, err := s.Login(context.Background(), "unknown", "bad")
	assert.Equal(t, errors.UnauthorizedError(""), err)
	token, err := s.Login(context.Background(), "user999", "pass")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func Test_service_authenticate(t *testing.T) {
	logger, _ := log.NewForTest()
	s := service{"test", 100, logger}
	assert.Nil(t, s.authenticate(context.Background(), "unknown", "bad"))
	assert.Nil(t, s.authenticate(context.Background(), "user", "pass"))
	assert.Nil(t, s.authenticate(context.Background(), "user10001", "pass"))
	assert.NotNil(t, s.authenticate(context.Background(), "user10", "pass"))
}

func Test_service_GenerateJWT(t *testing.T) {
	logger, _ := log.NewForTest()
	s := service{"test", 100, logger}
	token, err := s.generateJWT(entity.User{
		ID:   "100",
		Name: "user415",
	})
	if assert.Nil(t, err) {
		assert.NotEmpty(t, token)
	}
}
