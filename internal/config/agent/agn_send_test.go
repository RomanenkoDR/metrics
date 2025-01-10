package agent

import (
	"github.com/RomanenkoDR/metrics/internal/config/agent/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидался метод POST, но получен %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Ожидался Content-Type application/json, но получен %s", r.Header.Get("Content-Type"))
		}

		w.WriteHeader(http.StatusOK)
	}))
}

func TestSendReport(t *testing.T) {
	server := createTestServer(t) // создайте тестовый сервер
	defer server.Close()

	metrics := types.Metrics{
		ID:    "test",
		MType: "gauge",
		Value: 123.45,
	}

	err := sendReport(server.URL, metrics)
	if err != nil {
		t.Fatalf("Ошибка отправки отчета: %v", err)
	}
}
