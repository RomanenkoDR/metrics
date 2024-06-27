package main

import (
	"log"
	"time"

	config "github.com/RomanenkoDR/metrics/internal/config/agentConfig"
	logging "github.com/RomanenkoDR/metrics/internal/logging"
	metrics "github.com/RomanenkoDR/metrics/internal/metrics"
)

func main() {
	// Инициализация конфигурации агента.
	configuration := config.NewAgentConfig()
	configuration.InitAgentConfiguration()

	// Настройка логгера zap
	logger, err := logging.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	// Создание экземпляра метрик и агента
	metric := metrics.NewMetrics()
	agent := metrics.NewAgent(metric, time.Duration(configuration.ReportInterval), time.Duration(configuration.PollInterval), logger)

	// Запуск агента
	agent.Start()

	log.Printf("Агент запущен, начинаем сбор метрик для отправки на сервер %s...\n", configuration.Address)

	select {}
}
