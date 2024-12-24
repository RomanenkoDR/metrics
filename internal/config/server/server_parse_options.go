package server

import (
	"flag"
	"github.com/RomanenkoDR/metrics/internal/models"

	"github.com/caarlos0/env"
)

func ParseOptionsServer() (models.ConfigServer, error) {
	var cfg models.ConfigServer

	flag.StringVar(&cfg.Address,
		"a", "localhost:8080",
		"Add address and port in format <address>:<port>")
	flag.IntVar(&cfg.Interval,
		"i", 300,
		"Saving metrics to file interval")
	flag.StringVar(&cfg.Filename,
		"f", "/tmp/metrics-db.json",
		"File path")
	flag.BoolVar(&cfg.Restore,
		"r", true,
		"Restore metrics value from file")
	flag.StringVar(&cfg.DBDSN,
		"d",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		"Connection string in Postgres format")
	flag.StringVar(&cfg.Key, "k", "", "Sing key")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
