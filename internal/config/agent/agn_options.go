package agent

import "github.com/RomanenkoDR/metrics/internal/storage"

type Options struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
}

type Metrics struct {
	ID    string          `json:"id"`    // имя метрики
	MType string          `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta storage.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value storage.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}

var Encrypt bool
var Key []byte
