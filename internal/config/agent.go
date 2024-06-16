package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// AgentConfig содержит конфигурационные параметры для агента
type AgentConfig struct {
	Address        string        // Адрес сервера
	ReportInterval time.Duration // Интервал отправки метрик
	PollInterval   time.Duration // Интервал опроса метрик
}

// NewAgentConfig создает экземпляр AgentConfig с параметрами по умолчанию
func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		Address:        "localhost:8080", // Адрес сервера
		ReportInterval: 10 * time.Second, // Интервал отправки метрик
		PollInterval:   2 * time.Second,  // Интервал опроса метрик
	}
}

// Инициализирует конфигурацию агента, проверяя флаги командной строки и переменные окружения
func (c *AgentConfig) Init() {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&c.Address, "a", c.Address, "Адрес HTTP сервера")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "Интервал отправки метрик")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "Интервал опроса метрик")

	// Парсинг флагов командной строки
	flag.Parse()

	// Проверка наличия переменной окружения для адреса сервера
	if addr := os.Getenv("ADDRESS"); addr != "" {
		c.Address = addr
	}

	// Проверка наличия переменной окружения для интервала отправки метрик
	if reportIntervalStr := os.Getenv("REPORT_INTERVAL"); reportIntervalStr != "" {
		// Преобразовываем строковое значение в тип time.Duration и сохраняем результатт в переменную dur
		if dur, err := time.ParseDuration(reportIntervalStr); err == nil {
			c.ReportInterval = dur
		} else {
			// Вывод сообщения об ошибке в случае некорректного значения
			fmt.Printf("Некорректное значение REPORT_INTERVAL: %v\n", err)
		}
	}

	// Проверка наличия переменной окружения для интервала опроса метрик
	if pollIntervalStr := os.Getenv("POLL_INTERVAL"); pollIntervalStr != "" {
		// Преобразовываем строковое значение в тип time.Duration и сохраняем результатт в переменную dur
		if dur, err := time.ParseDuration(pollIntervalStr); err == nil {
			c.PollInterval = dur
		} else {
			// Вывод сообщения об ошибке в случае некорректного значения
			fmt.Printf("Некорректное значение POLL_INTERVAL: %v\n", err)
		}
	}
}
