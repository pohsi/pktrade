package main

import (
	"flag"
	"os"

	_ "github.com/lib/pq"
	"github.com/pohsi/pktrade/internal/config"
	"github.com/pohsi/pktrade/internal/server"
	"github.com/pohsi/pktrade/pkg/log"
)

const Version = "1.0.0"

func main() {

	l := log.New().With(nil, "version", Version)
	flagConfig := flag.String("config", "configs/dev.yml", "path to the config file")

	flag.Parse()
	cfg, err := config.Load(*flagConfig, l)
	if err != nil {
		l.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	s, err := server.New(*cfg, l, Version)
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
