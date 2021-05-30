package config

import (
	"io/ioutil"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pohsi/pktrade/pkg/env"
	"github.com/pohsi/pktrade/pkg/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultServerPort         = 8080 // Server port, defaults to 8080
	defaultJWTExpirationHours = 72   // JWT expiration in hours, defaults to 72 hours (3 days)
)

type Config struct {
	ServerPort int `yaml:"server_port" env:"SERVER_PORT"`

	DSN string `yaml:"dsn" env:"DSN,secret"`

	JWTSigningKey string `yaml:"jwt_signing_key" env:"JWT_SIGNING_KEY,secret"`

	JWTExpiration int `yaml:"jwt_expiration" env:"JWT_EXPIRATION"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.JWTSigningKey, validation.Required),
	)
}

func Load(file string, logger log.Logger) (*Config, error) {
	c := Config{
		ServerPort:    defaultServerPort,
		JWTExpiration: defaultJWTExpirationHours,
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	if err = env.New("APP_", logger.Infof).Load(&c); err != nil {
		return nil, err
	}

	if err = c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
