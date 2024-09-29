package agent

import (
	"flag"
	"github.com/caarlos0/env"
)

type Options struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func ParseOptions() (Options, error) {
	var opt Options

	flag.IntVar(&opt.PollInterval, "p", 2,
		"Frequensy in seconds for collecting metrics")
	flag.IntVar(&opt.ReportInterval, "r", 10,
		"Frequensy in seconds for sending report to the server")
	flag.StringVar(&opt.ServerAddress, "a", "localhost:8080",
		"Address of the server to send metrics")
	flag.Parse()

	err := env.Parse(&opt)
	if err != nil {
		return opt, err
	}

	return opt, nil
}
