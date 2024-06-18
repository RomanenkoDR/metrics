package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// DurationInSeconds - тип для интервала в секундах
type DurationInSeconds time.Duration

// String - метод для флага DurationInSeconds
func (d *DurationInSeconds) String() string {
	return fmt.Sprintf("%d", time.Duration(*d)/time.Second)
}

// Set - метод для флага DurationInSeconds
func (d *DurationInSeconds) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*d = DurationInSeconds(time.Duration(v) * time.Second)
	return nil
}

// AgentConfig содержит конфигурационные параметры для агента
type AgentConfig struct {
	Address        string            // Адрес сервера
	ReportInterval DurationInSeconds // Интервал отправки метрик
	PollInterval   DurationInSeconds // Интервал опроса метрик
}

// NewAgentConfig создает экземпляр AgentConfig с параметрами по умолчанию
func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		Address:        "localhost:8080",                    // Адрес сервера
		ReportInterval: DurationInSeconds(10 * time.Second), // Интервал отправки метрик
		PollInterval:   DurationInSeconds(2 * time.Second),  // Интервал опроса метрик
	}
}

// Инициализирует конфигурацию агента, проверяя флаги командной строки и переменные окружения
func (c *AgentConfig) InitAgentConfiguration() {
	// Определение флагов командной строки для настройки конфигурации.
	flag.StringVar(&c.Address, "a", c.Address, "Адрес HTTP сервера")
	flag.Var(&c.ReportInterval, "r", "Интервал отправки метрик в секундах")
	flag.Var(&c.PollInterval, "p", "Интервал опроса метрик в секундах")

	// Парсинг флагов командной строки
	flag.Parse()

	// Проверка наличия переменной окружения для адреса сервера
	if addr := os.Getenv("ADDRESS"); addr != "" {
		c.Address = addr
	}

	// Проверка наличия переменной окружения для интервала отправки метрик
	if reportIntervalStr := os.Getenv("REPORT_INTERVAL"); reportIntervalStr != "" {
		if dur, err := strconv.Atoi(reportIntervalStr); err == nil {
			c.ReportInterval = DurationInSeconds(time.Duration(dur) * time.Second)
		} else {
			fmt.Printf("Некорректное значение REPORT_INTERVAL: %v\n", err)
		}
	}

	// Проверка наличия переменной окружения для интервала опроса метрик
	if pollIntervalStr := os.Getenv("POLL_INTERVAL"); pollIntervalStr != "" {
		if dur, err := strconv.Atoi(pollIntervalStr); err == nil {
			c.PollInterval = DurationInSeconds(time.Duration(dur) * time.Second)
		} else {
			fmt.Printf("Некорректное значение POLL_INTERVAL: %v\n", err)
		}
	}
}
