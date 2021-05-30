// Package server provides running server and handle API request
package server

import (
	"github.com/pohsi/pktrade/pkg/log"
)

// Server is the major be responsible for run and handle rest request
type Server interface {
	Run() error

	Port() int
}

type concreteServer struct {
	logger log.Logger
	config Config
}

type Config struct {
	Port   int
	Logger log.Logger
}

// New creates server instance wich takes custom logger
func New(cfg Config) (Server, error) {

	if cfg.Logger == nil {
		cfg.Logger = log.New()
	}
	return &concreteServer{logger: cfg.Logger, config: cfg}, nil
}

func (c *concreteServer) Run() error {
	return nil
}

func (c *concreteServer) Port() int {
	return c.config.Port
}
