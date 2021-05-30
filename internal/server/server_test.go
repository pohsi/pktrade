package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewByConfig(t *testing.T) {
	port := 9999
	cfg := Config{Port: port}
	s, _ := New(cfg)
	assert.NotNil(t, s)
	assert.Equal(t, s.Port(), port)
}
