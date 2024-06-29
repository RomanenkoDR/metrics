package serverConfig

import (
	"flag"
	"os"
)

type Options struct {
	Address string
}

func ParseOptions() Options {
	var cfg Options

	flag.StringVar(&cfg.Address,
		"a", "localhost:8080",
		"Add addres and port in format <address>:<port>")
	flag.Parse()

	// get env vars
	if a := os.Getenv("ADDRESS"); a != "" {
		cfg.Address = a
	}

	return cfg
}
