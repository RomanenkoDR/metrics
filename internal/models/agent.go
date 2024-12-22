package models

import "github.com/RomanenkoDR/metrics/internal/storage"

type OptionsAgent struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Key            string `env:"KEY"`
}

type MetricsAgent struct {
	ID    string          `json:"id"`    // имя метрики
	MType string          `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta storage.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value storage.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}
