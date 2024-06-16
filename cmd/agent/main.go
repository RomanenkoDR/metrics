package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/RomanenkoDR/metrics/internal/config"
	metricagent "github.com/RomanenkoDR/metrics/internal/metricagent"
	"github.com/RomanenkoDR/metrics/internal/metrics"
)

func main() {
	// Инициализация конфигурации агента.
	cfg := config.NewAgentConfig()
	cfg.InitAgentConfiguration()

	// Изменение конфигурации на основе переменных окружения
	// // Проверяем наличие переменной окружения ADDRESS
	if addr := os.Getenv("ADDRESS"); addr != "" {
		cfg.Address = addr
	}

	// Проверяем наличие переменной окружения REPORT_INTERVAL
	if reportIntervalStr := os.Getenv("REPORT_INTERVAL"); reportIntervalStr != "" {

		// Если переменная окружения существует, пробуем преобразовать её значение в число
		if dur, err := strconv.Atoi(reportIntervalStr); err == nil {
			// Присваиваем значение в виде временного интервала
			cfg.ReportInterval = time.Duration(dur) * time.Second
		} else {
			// Логируем ошибку, если значение неверное
			log.Printf("Invalid value for REPORT_INTERVAL: %v", err)
		}
	}

	// Проверяем наличие переменной окружения POLL_INTERVAL
	if pollIntervalStr := os.Getenv("POLL_INTERVAL"); pollIntervalStr != "" {

		// Если переменная окружения существует, пробуем преобразовать её значение в число (секунды)
		if dur, err := strconv.Atoi(pollIntervalStr); err == nil {
			// Присваиваем значение в виде временного интервала
			cfg.PollInterval = time.Duration(dur) * time.Second
		} else {
			// Логируем ошибку, если значение неверное
			log.Printf("Invalid value for POLL_INTERVAL: %v", err)
		}
	}

	// Создание экземпляра метрик и агента
	metrics := metrics.NewMetrics()

	// Создание экземпляра агента с передачей экземпляра метрик и интервалов времени
	agent := metricagent.NewAgent(metrics, cfg.ReportInterval, cfg.PollInterval)

	// Запуск агента
	agent.Start()

	// Логирование сообщения о запуске агента и указание адреса сервера, на который отправляются метрики
	log.Printf("Агент запущен, собирает и отправляет метрики на сервер %s...\n", cfg.Address)

	// Блокируем основную горутину
	// Это позволяет приложению продолжать работу и ожидать внешних событий,
	// таких как сигналы операционной системы (например, завершение работы приложения) или другие события.
	// Блокировка основной горутины в этом месте обеспечивает бесконечный цикл работы агента и приложения в целом
	select {}
}
