package agent

import (
	"flag"
	"github.com/RomanenkoDR/metrics/internal/models"
	"github.com/caarlos0/env"
	"strings"
)

func ParseOptionsAgent() (models.ConfigAgent, error) {
	var cfg models.ConfigAgent
	cfg.Encrypt = false

	flag.IntVar(&cfg.PollInterval, "p", 2,
		"Frequensy in seconds for collecting metrics")
	flag.IntVar(&cfg.ReportInterval, "r", 10,
		"Frequensy in seconds for sending report to the server")
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080",
		"Address of the server to send metrics")
	flag.StringVar(&cfg.Key, "k", "",
		"Encryption key")
	flag.IntVar(&cfg.RateLimit, "l", 3,
		"Rate Limit")
	flag.Parse()

	cfg.ServerAddress = strings.Join([]string{"http:/", cfg.ServerAddress, "updates/"}, "/")

	if cfg.Key != "" {
		cfg.Encrypt = true
		cfg.KeyByte = []byte(cfg.Key)
	}

	err := env.Parse(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
