package servercfg

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Options struct {
	Address  string `env:"ADDRESS"`
	Interval int    `env:"STORE_INTERVAL"`
	Filename string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
	DBDSN    string `env:"DATABASE_DSN"`
}

// DB connection string parameters
const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func ParseOptions() (Options, error) {
	var cfg Options

	flag.StringVar(&cfg.Address,
		"a", "localhost:8080",
		"Add addres and port in format <address>:<port>")
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
		"d", fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s target_session_attrs=read-write",
			host, port, dbname, user, password),
		"Connection string in Postgres format")
	flag.Parse()

	// get env vars
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
