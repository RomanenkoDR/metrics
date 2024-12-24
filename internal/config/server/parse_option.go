package server

import (
	"flag"

	"github.com/caarlos0/env"
)

type Options struct {
	Address  string `env:"ADDRESS" envDefault:"localhost:8080"`
	Interval int    `env:"STORE_INTERVAL" envDefault:"300"`
	Filename string `env:"FILE_STORAGE_PATH" envDefault:"/cmd/internal/storage/metrics-db.json"`
	Restore  bool   `env:"RESTORE" envDefault:"true"`
	DBDSN    string `env:"DATABASE_DSN"`
	Key      string `env:"KEY"`
}

func ParseOptions() (Options, error) {
	var cfg Options

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
		"d", //fmt.Sprintf(
		//"host=%s port=%d dbname=%s user=%s password=%s target_session_attrs=read-write",
		//host, port, dbname, user, password),
		"",
		"Connection string in Postgres format")
	flag.StringVar(&cfg.Key, "k", "", "Sing key")
	flag.Parse()

	// get env vars
	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
