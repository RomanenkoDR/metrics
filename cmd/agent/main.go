// cmd/agent/main.go

package main

import (
	"log"

	"github.com/RomanenkoDR/metrics/internal/config"
	metricagent "github.com/RomanenkoDR/metrics/internal/metricagent"
	"github.com/RomanenkoDR/metrics/internal/metrics"
)

func main() {
	cfg := config.NewAgentConfig()
	cfg.Init()

	metrics := metrics.NewMetrics()
	agent := metricagent.NewAgent(metrics, cfg.ReportInterval, cfg.PollInterval)
	agent.Start()

	log.Println("Агент запущен, собирает и отправляет метрики...")

	select {} // Блокируем основную горутину
}
