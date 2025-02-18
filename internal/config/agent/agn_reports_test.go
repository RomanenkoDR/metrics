package agent

import (
   "context"
   "net/http"
   "net/http/httptest"
   "strings"
   "testing"

   "github.com/RomanenkoDR/metrics/internal/storage"
   "github.com/stretchr/testify/assert"
)

// Моковый HTTP-сервер
func mockServer(t *testing.T, expectedStatus int, expectedBody string) *httptest.Server {
   return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  assert.Equal(t, "POST", r.Method)
	  w.WriteHeader(expectedStatus)
	  w.Write([]byte(expectedBody))
   }))
}

// Тест sendRequest
func TestSendRequest(t *testing.T) {
   server := mockServer(t, http.StatusOK, `{"status":"ok"}`)
   defer server.Close()

   data := map[string]string{"test": "value"}
   err := sendRequest(server.URL, data)
   assert.NoError(t, err)
}

// Тест sendReport
func TestSendReport(t *testing.T) {
   server := mockServer(t, http.StatusOK, `{"status":"ok"}`)
   defer server.Close()

   metric := Metrics{ID: "test_metric", MType: counterType, Delta: 42}
   err := sendReport(server.URL, metric)
   assert.NoError(t, err)
}

// Тест sendReportBatch
func TestSendReportBatch(t *testing.T) {
   server := mockServer(t, http.StatusOK, `{"status":"ok"}`)
   defer server.Close()

   metrics := []Metrics{
	  {ID: "metric1", MType: counterType, Delta: 10},
	  {ID: "metric2", MType: gaugeType, Value: 5.5},
   }
   err := sendReportBatch(server.URL, metrics)
   assert.NoError(t, err)
}

// Тест ProcessReport
func TestProcessReport(t *testing.T) {
   server := mockServer(t, http.StatusOK, `{"status":"ok"}`)
   defer server.Close()

   memStore := storage.MemStorage{
	  CounterData: map[string]storage.Counter{"counter1": 100}, // Используем storage.Counter как int64
	  GaugeData:   map[string]storage.Gauge{"gauge1": 42.42},   // Используем storage.Gauge как float64
   }

   err := ProcessReport(strings.TrimPrefix(server.URL, "http://"), memStore)
   assert.NoError(t, err)
}

// Тест ProcessBatch
func TestProcessBatch(t *testing.T) {
   server := mockServer(t, http.StatusOK, `{"status":"ok"}`)
   defer server.Close()

   memStore := storage.MemStorage{
	  CounterData: map[string]storage.Counter{"counter1": 100}, // Используем storage.Counter как int64
	  GaugeData:   map[string]storage.Gauge{"gauge1": 42.42},   // Используем storage.Gauge как float64
   }

   err := ProcessBatch(context.Background(), strings.TrimPrefix(server.URL, "http://"), memStore)
   assert.NoError(t, err)
}
