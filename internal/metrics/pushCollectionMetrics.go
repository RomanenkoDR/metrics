package metrics

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	metrics                     *SystemMetrics
	updateMetricCounterInterval time.Duration // Интервал опроса
	pushMetricInterval          time.Duration // Интервал отправки
}

// Создаем новый экземпляр агента с указанными интервалами опроса и отправки метрик
func NewAgent(m *SystemMetrics, updateMetricCounterInterval, pushMetricInterval time.Duration) *Agent {
	return &Agent{
		metrics:                     m,
		updateMetricCounterInterval: updateMetricCounterInterval,
		pushMetricInterval:          pushMetricInterval,
	}
}

// Запуск агента, который регулярно собирает и отправляет метрики
func (a *Agent) Start() {

	tickerUpdateMetric := time.NewTicker(a.updateMetricCounterInterval)
	tickerPushMetric := time.NewTicker(a.pushMetricInterval)

	// Запускает асинхронную горутину, которая будет работать параллельно основной программе
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

// Собирает все метрики и отправляет их на сервер
func (a *Agent) CollectionAllMetrics() {
	metrics := a.metrics.CopyMetricsBeforePush()
	for name, metric := range metrics {
		a.pushMetricsToServer(name, metric)
	}
}

// Создаем копию метрик для отправки, что бы внешние изменения не повлияли на оригинальные данные
func (m *SystemMetrics) CopyMetricsBeforePush() map[string]SystemMetric {
	// Создаем новую мапу с тем же размером, что и "m.metrics"
	copyMetrics := make(map[string]SystemMetric, len(m.metrics))
	for key, value := range m.metrics {
		copyMetrics[key] = value
	}
	// Возвращаем созданную копию метрик
	return copyMetrics
}

// Функция отправлки метрик на сервер
func (a *Agent) pushMetricsToServer(name string, metric SystemMetric) {

	var data string

	// Формируем URL для отправки метрики в зависимости от её типа
	if metric.Type == Gauge {
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", metric.Type, name, metric.Value.(float64))
	} else if metric.Type == Counter {
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", metric.Type, name, metric.Value.(int64))
	} else {
		log.Printf("Неизвестный тип метрики: %s", metric.Type)
		return
	}

	// // Выводим информацию о том, какая метрика отправляется
	// log.Printf("Отправка метрики: %s", data)

	// Создаем HTTP запрос
	req, err := http.NewRequest("POST", data, nil)
	if err != nil {
		// Выводим сообщение об ошибке при создании запроса
		log.Printf("Ошибка создания запроса: %v", err)
		return
	}
	// Устанавливаем заголовки в Header запроса
	req.Header.Set("Content-Type", "text/plain")

	// Отправляем запрос на сервер
	client := &http.Client{}
	// Отправляет HTTP запрос на сервер и получает ответ. Если возникает ошибка во время отправки запроса, она логируется.
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка отправки запроса: %v", err)
		return
	}

	// Отложенное закрытие тела ответа
	defer resp.Body.Close()

	// Проверяем статус ответа сервера
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка ответа сервера: %v", resp.StatusCode)
	}
}
