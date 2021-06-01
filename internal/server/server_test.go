package server

import (
	"testing"

	"github.com/pohsi/pktrade/internal/config"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestNewByConfig(t *testing.T) {
	port := 9999
	cfg := config.Config{ServerPort: port}
	logger, _ := log.NewForTest()
	s, _ := New(cfg, logger, `test`)
	assert.NotNil(t, s)
	assert.Equal(t, s.Port(), port)
}
