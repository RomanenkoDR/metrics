package main

import (
	"log"
	"os"
	"strconv"
	"time"

	config "github.com/RomanenkoDR/metrics/internal/config/agentConfig"
	metricagent "github.com/RomanenkoDR/metrics/internal/metricagent"
	"github.com/RomanenkoDR/metrics/internal/metrics"
)

func main() {
	// Инициализация конфигурации агента.
	configuration := config.NewAgentConfig()
	configuration.InitAgentConfiguration()

	// Проверка переменных окружения и их установка
	if addr := os.Getenv("ADDRESS"); addr != "" {
		configuration.Address = addr
	}

	if reportIntervalStr := os.Getenv("REPORT_INTERVAL"); reportIntervalStr != "" {
		if dur, err := strconv.Atoi(reportIntervalStr); err == nil {
			configuration.ReportInterval = config.DurationInSeconds(time.Duration(dur) * time.Second)
		} else {
			log.Printf("Некорректно задан интервал REPORT_INTERVAL: %v", err)
		}
	}

	if pollIntervalStr := os.Getenv("POLL_INTERVAL"); pollIntervalStr != "" {
		if dur, err := strconv.Atoi(pollIntervalStr); err == nil {
			configuration.PollInterval = config.DurationInSeconds(time.Duration(dur) * time.Second)
		} else {
			log.Printf("Некорректно задан интервал POLL_INTERVAL: %v", err)
		}
	}

	// Создание экземпляра метрик и агента
	metrics := metrics.NewMetrics()
	agent := metricagent.NewAgent(metrics, time.Duration(configuration.ReportInterval), time.Duration(configuration.PollInterval))

	// Запуск агента
	agent.Start()

	log.Printf("Агент запущен, начинаем сбор метрик для отправки на сервер %s...\n", configuration.Address)

	select {}
}
