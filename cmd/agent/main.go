package main

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/config/agent"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

type Sender func(context.Context, string, storage.MemStorage) error

func Retry(sender Sender, retries int, delay time.Duration) Sender {
	return func(ctx context.Context, serverAddress string, m storage.MemStorage) error {
		for r := 0; ; r++ {
			err := sender(ctx, serverAddress, m)
			if err == nil || r >= retries {
				logger.DebugLogger.Sugar().Infof("Completed retry cycle with %d attempts", r)
				return err
			}

			logger.DebugLogger.Sugar().Warnf("Function call failed: %v, retrying in %v", err, delay)

			delay += time.Second * 2

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				logger.DebugLogger.Sugar().Error("Context cancelled, stopping retries")
				return ctx.Err()
			}
		}
	}
}

func main() {
	logger.DebugLogger.Sugar().Info("Starting the application")

	// Parse cli options
	cfg, err := agent.ParseOptions()
	if err != nil {
		logger.DebugLogger.Sugar().Fatal("Failed to parse options: ", err)
	}

	// Initiate tickers
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	// Initiate new storage
	m := storage.New()

	logger.DebugLogger.Sugar().Info("Initialization completed, starting main loop")

	// Collect data from MemStats and send to the server
	for {
		select {
		case <-pollTicker.C:
			logger.DebugLogger.Sugar().Debug("Polling memory stats")
			agent.ReadMemStats(&m)
		case <-reportTicker.C:
			logger.DebugLogger.Sugar().Debug("Reporting memory stats")
			fn := Retry(agent.ProcessBatch, 3, 1*time.Second)
			err := fn(context.Background(), cfg.ServerAddress, m)
			if err != nil {
				logger.DebugLogger.Sugar().Error("Failed to process batch: ", err)
			}
		}
	}
}
