package main

import (
	"log"
	"time"

	config "github.com/RomanenkoDR/metrics/internal/config/agentConfig"
	metrics "github.com/RomanenkoDR/metrics/internal/metrics"
)

func main() {
	// Инициализация конфигурации агента.
	configuration := config.NewAgentConfig()
	configuration.InitAgentConfiguration()

	// Создание экземпляра метрик и агента
	metric := metrics.NewMetrics()
	agent := metrics.NewAgent(metric, time.Duration(configuration.ReportInterval), time.Duration(configuration.PollInterval))

	// Запуск агента
	agent.Start()

	log.Printf("Агент запущен, начинаем сбор метрик для отправки на сервер %s...\n", configuration.Address)

	select {}
}
