package metricagent

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RomanenkoDR/metrics/internal/metrics"
)

type Agent struct {
	metrics                     *metrics.SystemMetrics
	updateMetricCounterInterval time.Duration // Интервал опроса
	pushMetricInterval          time.Duration // Интервал отправки
}

// Создаем новый экземпляр агента с указанными интервалами опроса и отправки метрик
func NewAgent(m *metrics.SystemMetrics, updateMetricCounterInterval, pushMetricInterval time.Duration) *Agent {
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
			case <-tickerUpdateMetric.C: // Выполняется когда срабатывает update
				// Логируем начало сбора метрик для информирования
				log.Println("Сбор метрик...")
				// Вызываем метод Collection() объекта метрик для сбора актуальных данных
				a.metrics.CollectionOfMetrics()
			case <-tickerPushMetric.C:
				log.Println("Отправка метрик на сервер...")
				// Вызываем метод CollectionAllMetrics() для отправки всех собранных метрик на сервер
				a.CollectionAllMetrics() // Отправка метрик на сервер
			}
		}
	}()
}

// Собирает все метрики и отправляет их на сервер
func (a *Agent) CollectionAllMetrics() {
	// Собираем все метрики
	metrics := a.metrics.GetMetrics()
	// Перебираем все метрики
	for name, metric := range metrics {
		// Отправляем каждую метрику на сервер
		a.pushMetricsToServer(name, metric)
	}
}

// Функция отправлки метрик на сервер
func (a *Agent) pushMetricsToServer(name string, metric metrics.SystemMetric) {

	var data string

	// Формируем URL для отправки метрики в зависимости от её типа
	if metric.Type == metrics.Gauge { // Если у метрики тип Gauge
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", metric.Type, name, metric.Value.(float64))
	} else if metric.Type == metrics.Counter { // Если у метрики тип Counter
		data = fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", metric.Type, name, metric.Value.(int64))
	} else {
		// Выводим сообщение об ошибке в случае неизвестного типа метрики
		log.Printf("Неизвестный тип метрики: %s", metric.Type)
		return
	}

	// Выводим информацию о том, какая метрика отправляется
	log.Printf("Отправка метрики: %s", data)

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
		// Выводим сообщение об ошибке при отправке запроса
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
