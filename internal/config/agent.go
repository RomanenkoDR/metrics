// internal/config/agent.go

package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type AgentConfig struct {
	Address        string        // Адрес сервера
	ReportInterval time.Duration // Интервал отправки метрик
	PollInterval   time.Duration // Интервал опроса метрик
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
	}
}

func (c *AgentConfig) Init() {
	flag.StringVar(&c.Address, "a", c.Address, "address of the HTTP server")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "report interval")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "poll interval")

	flag.Parse()

	if addr := os.Getenv("ADDRESS"); addr != "" {
		c.Address = addr
	}

	if reportIntervalStr := os.Getenv("REPORT_INTERVAL"); reportIntervalStr != "" {
		if dur, err := time.ParseDuration(reportIntervalStr); err == nil {
			c.ReportInterval = dur
		} else {
			fmt.Printf("Invalid value for REPORT_INTERVAL: %v\n", err)
		}
	}

	if pollIntervalStr := os.Getenv("POLL_INTERVAL"); pollIntervalStr != "" {
		if dur, err := time.ParseDuration(pollIntervalStr); err == nil {
			c.PollInterval = dur
		} else {
			fmt.Printf("Invalid value for POLL_INTERVAL: %v\n", err)
		}
	}
}
