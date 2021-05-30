package main

import (
	"flag"
	"os"

	"github.com/pohsi/pktrade/internal/config"
	"github.com/pohsi/pktrade/internal/server"
	"github.com/pohsi/pktrade/pkg/log"
)

type Result int

func main() {

	l := log.New()
	flagConfig := flag.String("config", "configs/dev.yml", "path to the config file")

	flag.Parse()
	cfg, err := config.Load(*flagConfig, l)
	if err != nil {
		l.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	s, err := server.New(server.Config{Port: cfg.ServerPort, Logger: l})
	if err != nil {
		l.Error(err)
		os.Exit(-1)
	}
	err = s.Run()
	if err != nil {
		l.Error(err)
		os.Exit(-1)
	}
}
