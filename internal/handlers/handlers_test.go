package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	memstorage "github.com/RomanenkoDR/metrics/internal/storage/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationUpdateMetric(t *testing.T) {
	// Создаем реальное хранилище для метрик
	storage := memstorage.NewMemStorage()

	// Имитируем ваш HTTP сервер с помощью Chi
	r := chi.NewRouter()
	r.HandleFunc("/metrics/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		UpdateMetric(w, r, storage) // Обработчик, который будет вызываться при запросах на /metrics/{metricType}/{metricName}/{metricValue}
	})

	// Запускаем ваш HTTP сервер в тестовом режиме
	server := httptest.NewServer(r) // Создание тестового HTTP сервера с маршрутами, определенными в роутере `r`
	defer server.Close()            // Закрытие сервера после завершения теста

	// Определяем тестовые случаи
	tests := []struct {
		name         string // Название теста
		url          string // URL для отправки запроса
		expectedCode int    // Ожидаемый HTTP код ответа
	}{
		{
			name:         "Проверка отправки метрики Gauge",
			url:          fmt.Sprintf("%s/metrics/gauge/my_metric_name/123.45", server.URL), // URL для отправки метрики типа gauge
			expectedCode: http.StatusOK,                                                     // Ожидаемый код ответа HTTP 200 OK
		},
		{
			name:         "Проверка отправки метрики Counter",
			url:          fmt.Sprintf("%s/metrics/counter/my_counter/5", server.URL), // URL для отправки метрики типа counter
			expectedCode: http.StatusOK,                                              // Ожидаемый код ответа HTTP 200 OK
		},
		{
			name:         "Проверка некорректного типа метрики",
			url:          fmt.Sprintf("%s/metrics/invalid_type/my_metric_name/123.45", server.URL), // URL с некорректным типом метрики
			expectedCode: http.StatusBadRequest,                                                    // Ожидаемый код ответа HTTP 400 Bad Request
		},
		{
			name:         "Проверка пустой метрики без значения",
			url:          fmt.Sprintf("%s/metrics/gauge/my_metric_name/invalid_value", server.URL), // URL с некорректным значением метрики типа gauge
			expectedCode: http.StatusBadRequest,                                                    // Ожидаемый код ответа HTTP 400 Bad Request
		},
	}

	// Проходим по всем тестовым случаям
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, test.url, nil) // Создание нового HTTP запроса POST по указанному URL
			assert.NoError(t, err)                                      // Проверка отсутствия ошибок при создании запроса

			res, err := http.DefaultClient.Do(req) // Выполнение HTTP запроса с использованием клиента по умолчанию
			assert.NoError(t, err)                 // Проверка отсутствия ошибок при выполнении запроса
			defer res.Body.Close()                 // Закрытие тела ответа после завершения проверок

			// Проверяем код ответа
			assert.Equal(t, test.expectedCode, res.StatusCode)

			// Проверяем, что метрика успешно добавлена в хранилище (для успешных случаев)
			if test.expectedCode == http.StatusOK {
				switch {
				case test.url[:15] == fmt.Sprintf("%s/metrics/gauge/", server.URL):
					value := storage.GetGauge("my_metric_name")                               // Получаем текущее значение метрики типа gauge из хранилища
					expectedValue, _ := strconv.ParseFloat(test.url[len(server.URL)+23:], 64) // Ожидаемое значение метрики, извлеченное из URL
					assert.InDelta(t, expectedValue, value, 0.001)                            // Проверка на равенство с заданной точностью (допускается погрешность 0.001)
				case test.url[:18] == fmt.Sprintf("%s/metrics/counter/", server.URL):
					value := storage.GetCounter("my_counter")                                   // Получаем текущее значение метрики типа counter из хранилища
					expectedValue, _ := strconv.ParseInt(test.url[len(server.URL)+25:], 10, 64) // Ожидаемое значение метрики, извлеченное из URL
					assert.Equal(t, expectedValue, value)                                       // Проверка на точное равенство ожидаемого и полученного значения
				}
			}
		})
	}
}
