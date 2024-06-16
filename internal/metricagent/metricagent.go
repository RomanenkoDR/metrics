package agent

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RomanenkoDR/metrics/internal/metrics"
)

type Agent struct {
	metrics                     *metrics.Metrics
	updateMetricCounterInterval time.Duration // Интервал опроса
	pushMetricInterval          time.Duration // Интервал отправки
}

func NewAgent(m *metrics.Metrics, updateMetricCounterInterval, pushMetricInterval time.Duration) *Agent {
	return &Agent{
		metrics:                     m,
		updateMetricCounterInterval: updateMetricCounterInterval,
		pushMetricInterval:          pushMetricInterval,
	}
}

func (a *Agent) Start() {
	tickerUpdateMetric := time.NewTicker(a.updateMetricCounterInterval)
	tickerPushMetric := time.NewTicker(a.pushMetricInterval)

	go func() {
		for {
			select {
			case <-tickerUpdateMetric.C: // Выполняется когда срабатывает update
				log.Println("Сбор метрик...")
				a.metrics.Collection() // Сбор метрик
			case <-tickerPushMetric.C:
				log.Println("Отправка метрик на сервер...")
				a.CollectionAllMetrics() // Отправка метрик на сервер
			}
		}
	}()
}

func (a *Agent) CollectionAllMetrics() {
	metrics := a.metrics.GetMetrics()
	for name, metric := range metrics {
		a.pushMetricsToServer(name, metric)
	}
}

func (a *Agent) pushMetricsToServer(name string, metric metrics.Metric) {

	var data string

	if metric.Type == metrics.Gauge {
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", metric.Type, name, metric.Value.(float64))
	} else if metric.Type == metrics.Counter {
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", metric.Type, name, metric.Value.(int64))
	} else {
		log.Printf("Неизвестный тип метрики: %s", metric.Type)
		return
	}

	log.Printf("Отправка метрики: %s", data)

	req, err := http.NewRequest("POST", data, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса: %v", err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")

	// req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка отправки запроса: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка ответа сервера: %v", resp.StatusCode)
	}
}

func main() {
	metrics := metrics.NewMetrics()
	agent := NewAgent(metrics, 2*time.Second, 10*time.Second)

	agent.Start()
	log.Println("Соединение с сервером: http://localhost:8080")

	select {}

}
