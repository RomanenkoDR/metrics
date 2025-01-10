package server

import (
	"context"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"net/http"
	"testing"
	"time"
)

func TestSetupShutdown(t *testing.T) {

	server := &http.Server{}
	ctx, cancel := context.WithCancel(context.Background())
	store := &storage.Localfile{Path: "/tmp/test_metrics.json"}
	metrics := storage.New()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	setupShutdown(ctx, cancel, server, store, &metrics)
}
