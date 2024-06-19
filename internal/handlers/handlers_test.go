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
		// Обработчик, который будет вызываться при запросах на /metrics/{metricType}/{metricName}/{metricValue}
		UpdateMetric(w, r, storage)
	})

	// Запускаем ваш HTTP сервер в тестовом режиме
	// Создание тестового HTTP сервера с маршрутами, определенными в роутере `r`
	server := httptest.NewServer(r)
	// Закрытие сервера после завершения теста
	defer server.Close()

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
			// Создание нового HTTP запроса POST по указанному URL
			req, err := http.NewRequest(http.MethodPost, test.url, nil)
			// Проверка отсутствия ошибок при создании запроса
			assert.NoError(t, err)
			// Выполнение HTTP запроса с использованием клиента по умолчанию
			res, err := http.DefaultClient.Do(req)
			// Проверка отсутствия ошибок при выполнении запроса
			assert.NoError(t, err)
			// Закрытие тела ответа после завершения проверок
			defer res.Body.Close()
			// Проверяем код ответа
			assert.Equal(t, test.expectedCode, res.StatusCode)

			// Проверяем, что метрика успешно добавлена в хранилище (для успешных случаев)
			if test.expectedCode == http.StatusOK {
				switch {
				case test.url[:15] == fmt.Sprintf("%s/metrics/gauge/", server.URL):
					// Получаем текущее значение метрики типа gauge из хранилища
					value := storage.GetGauge("my_metric_name")
					// Ожидаемое значение метрики, извлеченное из URL
					expectedValue, _ := strconv.ParseFloat(test.url[len(server.URL)+23:], 64)
					// Проверка на равенство с заданной точностью (допускается погрешность 0.001)
					assert.InDelta(t, expectedValue, value, 0.001)
				case test.url[:18] == fmt.Sprintf("%s/metrics/counter/", server.URL):
					// Получаем текущее значение метрики типа counter из хранилища
					value := storage.GetCounter("my_counter")
					// Ожидаемое значение метрики, извлеченное из URL
					expectedValue, _ := strconv.ParseInt(test.url[len(server.URL)+25:], 10, 64)
					assert.Equal(t, expectedValue, value)
					// Проверка на точное равенство ожидаемого и полученного значения
				}

			}
		})
	}
}
