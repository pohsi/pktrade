package config

import (
	"testing"

	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCongif(t *testing.T) {
	cfg := Config{}
	assert.Equal(t, 0, cfg.ServerPort)
	assert.Equal(t, "", cfg.DSN)
	assert.Equal(t, "", cfg.JWTSigningKey)
	assert.Equal(t, 0, cfg.JWTExpiration)
}

func TestLoadAbsenceFile(t *testing.T) {
	cfg, err := Load("", log.New())
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

const (
	testdataPath = "../testdata/configs/"
)

func TestLoadGoodFile(t *testing.T) {

	cfg, err := Load(testdataPath+"test_config.yml", log.New())
	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 7777, cfg.ServerPort)
	assert.Equal(t, "postgres://127.0.0.1/pktrade?sslmode=disable&user=postgres&password=postgres", cfg.DSN)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", cfg.JWTSigningKey)
	assert.Equal(t, 24, cfg.JWTExpiration)
}

func TestLoadRequiredFieldMissingFile(t *testing.T) {
	cfg, err := Load(testdataPath+"test_config2.yml", log.New())
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}

func TestLoadNoneRequiredFieldMissingFile(t *testing.T) {
	cfg, err := Load(testdataPath+"test_config2.yml", log.New())
	assert.NotNil(t, err)
	assert.Nil(t, cfg)
}
