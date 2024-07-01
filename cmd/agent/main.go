package main

import (
	"log"
	"time"

	agentConfigPcg "github.com/RomanenkoDR/metrics/internal/config/agentConfig"
	metricsPcg "github.com/RomanenkoDR/metrics/internal/metrics"
	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

func main() {

	// Парсинг опций командной строки
	// Такие как интервалы опроса и отчета,
	// адрес сервера и другие настройки. Если происходит ошибка при парсинге,
	// программа завершает работу с помощью panic.
	config, err := agentConfigPcg.ParseOptions()
	if err != nil {
		panic(err)
	}

	// Инициализация тикеров
	// pollTicker и reportTicker используются для периодического запуска задач.
	// pollTicker запускает задачу по сбору данных
	// reportTicker запускает задачу по отправке собранных данных на сервер.
	// Интервалы задаются в конфигурации.
	pollTicker := time.NewTicker(time.Second * time.Duration(config.PollInterval))
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Second * time.Duration(config.ReportInterval))
	defer reportTicker.Stop()

	// Инициализация нового хранилища
	m := memStoragePcg.New()

	// Сбор данных из MemStats и отправка на сервер
	for {
		select {
		case <-pollTicker.C:
			metricsPcg.ReadMemStats(&m) // Сбор данных из MemStats
		case <-reportTicker.C:
			err := metricsPcg.PushMetricsToServer(config.ServerAddress, m) // Отправка данных на сервер
			if err != nil {
				log.Println(err)
			}
		}
	}
}
