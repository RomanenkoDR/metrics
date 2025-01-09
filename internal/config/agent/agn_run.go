package agent

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"sync"
	"time"
)

// RunAgent запускает выполнение агента.
// Основная функция запускает сбор и отправку метрик, управляет завершением работы через graceful shutdown.
func RunAgent() {
	// Парсим конфигурацию
	cfg, err := ParseOptions()
	if err != nil {
		log.Fatalf("Ошибка парсинга конфигурации: %v", err)
	}

	// Инициализация хранилища метрик
	m := storage.New()

	// Инициализация контекста для управления горутинами
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Каналы для передачи метрик
	metricsCh := make(chan storage.MemStorage, cfg.RateLimit)
	defer close(metricsCh)

	// WaitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// Запуск горутины для сбора метрик через runtime
	wg.Add(1)
	go func() {
		defer wg.Done()
		startCollecting(ctx, m, metricsCh, time.Second*time.Duration(cfg.PollInterval))
	}()

	// Запуск горутины для сбора системных метрик через gopsutil
	wg.Add(1)
	go func() {
		defer wg.Done()
		startSystemMetricsCollecting(ctx, m, time.Second*time.Duration(cfg.PollInterval))
	}()

	// Запуск горутин для отправки метрик на сервер
	for w := 1; w <= cfg.RateLimit; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startReporting(ctx, cfg, metricsCh, time.Second*time.Duration(cfg.ReportInterval))
		}()
	}

	// Ожидание сигнала завершения работы
	waitForShutdown(cancel)

	// Ожидание завершения всех горутин
	wg.Wait()
}
