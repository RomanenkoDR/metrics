package metrics

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// Agent represents the metrics agent
type Agent struct {
	metrics                     *SystemMetrics
	updateMetricCounterInterval time.Duration
	pushMetricInterval          time.Duration
	logger                      *zap.Logger
	serverAddress               string
}

// NewAgent creates a new agent
func NewAgent(m *SystemMetrics, updateMetricCounterInterval, pushMetricInterval time.Duration, logger *zap.Logger) *Agent {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "http://localhost:8080"
	}

	return &Agent{
		metrics:                     m,
		updateMetricCounterInterval: updateMetricCounterInterval,
		pushMetricInterval:          pushMetricInterval,
		logger:                      logger,
		serverAddress:               serverAddress,
	}
}

// Start runs the agent
func (a *Agent) Start() {

	tickerUpdateMetric := time.NewTicker(a.updateMetricCounterInterval)
	tickerPushMetric := time.NewTicker(a.pushMetricInterval)

	go func() {
		for {
			select {
			case <-tickerUpdateMetric.C:
				a.metrics.CollectionOfMetrics()
			case <-tickerPushMetric.C:
				a.CollectionAllMetrics()
			}
		}
	}()
}

// CollectionAllMetrics collects and pushes all metrics
func (a *Agent) CollectionAllMetrics() {
	metrics := a.metrics.CopyMetricsBeforePush()
	for name, metric := range metrics {
		a.pushMetricsToServer(name, metric)
	}
}

// CopyMetricsBeforePush copies metrics before pushing
func (m *SystemMetrics) CopyMetricsBeforePush() map[string]SystemMetric {
	copyMetrics := make(map[string]SystemMetric, len(m.metrics))
	for key, value := range m.metrics {
		copyMetrics[key] = value
	}
	return copyMetrics
}

// pushMetricsToServer pushes metrics to server
func (a *Agent) pushMetricsToServer(name string, metric SystemMetric) {
	client := resty.New()

	var url string

	if metric.Type == Gauge {
		url = fmt.Sprintf("%s/update/%s/%s/%f", a.serverAddress, metric.Type, name, metric.Value.(float64))
	} else if metric.Type == Counter {
		url = fmt.Sprintf("%s/update/%s/%s/%d", a.serverAddress, metric.Type, name, metric.Value.(int64))
	} else {
		a.logger.Error("Неизвестный тип метрики", zap.String("type", metric.Type))
		return
	}

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)
	if err != nil {
		a.logger.Error("Ошибка отправки запроса", zap.Error(err))
		return
	}

	if resp.StatusCode() != http.StatusOK {
		a.logger.Info("Метрика успешно отправлена", zap.String("name", name))
	}
}
