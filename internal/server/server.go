package server

import (
	"fmt"

	"github.com/pohsi/pktrade/pkg/log"
)

type Server struct {
	logger log.Logger
}

func New(logger log.Logger) *Server {
	return &Server{logger: logger}
}

func (s Server) Run() int {
	fmt.Println("server run")
	return 0
}
