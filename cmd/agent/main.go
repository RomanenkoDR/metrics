package main

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/agent"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfg agent.Options
	// Логируем старт приложения.
	cfg, err := agent.ParseOptions()
	if err != nil {
		panic(err)
	}

	// Initiate new storage
	m := storage.New()

	// Init channels
	done := make(chan struct{})
	metricsCh := make(chan storage.MemStorage, cfg.RateLimit)
	defer close(metricsCh)

	// Collect data from MemStats and send to the server
	// Gather facts
	go func(timer time.Duration) {
		for {
			time.Sleep(timer)
			agent.ReadMemStats(&m, metricsCh)
		}
	}(time.Second * time.Duration(cfg.PollInterval))

	// Send metrics to the server
	for w := 1; w <= cfg.RateLimit; w++ {
		go func(timer time.Duration) {
			for {
				time.Sleep(timer)
				fn := agent.Retry(agent.ProcessBatch, 3, 1*time.Second)
				err := fn(context.Background(), cfg, metricsCh)
				if err != nil {
					log.Println(err)
				}
			}
		}(time.Second * time.Duration(cfg.ReportInterval))
	}

	// Gracefull shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		close(done)
	}()

	<-done
}
