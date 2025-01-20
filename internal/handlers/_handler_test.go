package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тестирование обработчика HandleMain
func TestHandleMain(t *testing.T) {
	handler := NewHandler()

	// Мокаем данные, которые вернет Store
	handler.Store.GaugeData = map[string]storage.Gauge{
		"metric1": 12.34,
	}
	handler.Store.CounterData = map[string]storage.Counter{
		"metric2": 10,
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем записи для ответа
	rr := httptest.NewRecorder()
	handler.HandleMain(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Тестирование обработчика HandlePing
func TestHandlePing(t *testing.T) {
	handler := NewHandler()

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Мокаем функцию Ping базы данных
	handler.DBconn = nil // Можно заменить на мок-объект

	rr := httptest.NewRecorder()
	handler.HandlePing(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)
}

// Тестирование обработчика HandleUpdate
func TestHandleUpdate(t *testing.T) {
	handler := NewHandler()

	// Мокаем Store
	handler.Store.CounterData = make(map[string]storage.Counter)
	handler.Store.GaugeData = make(map[string]storage.Gauge)

	req, err := http.NewRequest("POST", "/update/counter/metric1/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUpdate(rr, req)

	// Проверяем статус код и обновление значения
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, storage.Counter(10), handler.Store.CounterData["metric1"])
}

// Тестирование обработчика HandleUpdateJSON
func TestHandleUpdateJSON(t *testing.T) {
	handler := NewHandler()

	// Мокаем данные JSON для теста
	metric := Metrics{
		ID:    "metric1",
		MType: counterType,
		Delta: new(int64),
	}
	*metric.Delta = 10
	reqBody, err := json.Marshal(metric)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/update/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleUpdateJSON(rr, req)

	// Проверяем статус код и обновление значения
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, storage.Counter(10), handler.Store.CounterData["metric1"])
}

// Тестирование обработчика HandleValue
func TestHandleValue(t *testing.T) {
	handler := NewHandler()

	// Мокаем данные для проверки
	handler.Store.GaugeData = map[string]storage.Gauge{
		"metric1": 12.34,
	}

	req, err := http.NewRequest("GET", "/value/metric1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleValue(rr, req)

	// Проверяем статус код и содержимое ответа
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "12.34", rr.Body.String())
}

// Тестирование обработчика HandleValueJSON
func TestHandleValueJSON(t *testing.T) {
	handler := NewHandler()

	// Мокаем данные для проверки
	handler.Store.GaugeData = map[string]storage.Gauge{
		"metric1": 12.34,
	}

	// Создаем запрос JSON
	metric := Metrics{
		ID:    "metric1",
		MType: gaugeType,
	}
	reqBody, err := json.Marshal(metric)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/value/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.HandleValueJSON(rr, req)

	// Проверяем статус код и содержимое ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp Metrics
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "metric1", resp.ID)
	assert.Equal(t, float64(12.34), *resp.Value)
}
