package main

import (
	"github.com/pohsi/pktrade/internal/server"
	"github.com/pohsi/pktrade/pkg/log"
)

func main() {
	s := server.New(log.New())
	s.Run()
}
