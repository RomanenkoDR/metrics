package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	s "github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestProcessReport(t *testing.T) {
	responseBody := "response"

	tests := []struct {
		name     string
		store    s.MemStorage
		wanterr  bool // Проверяем только факт наличия ошибки
		wantcode int
	}{
		{
			name: "Test Valid Post request gauge metric",
			store: s.MemStorage{
				GaugeData: map[string]s.Gauge{
					"valid": s.Gauge(2.32),
				},
			},
			wanterr:  false, // Ошибки не должно быть
			wantcode: http.StatusOK,
		},
		{
			name:     "Test Empty metric",
			store:    s.MemStorage{CounterData: map[string]s.Counter{}},
			wanterr:  true, // Ожидаем ошибку из-за отсутствия данных
			wantcode: http.StatusBadRequest,
		},
		{
			name: "Test Invalid Post request counter metric",
			store: s.MemStorage{
				CounterData: map[string]s.Counter{
					"valid": s.Counter(2),
				},
			},
			wanterr:  true, // Ожидаем ошибку
			wantcode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				http.Error(rw, responseBody, tc.wantcode)
			}))
			defer server.Close()

			err := ProcessReport(strings.Replace(server.URL, "http://", "", 1), tc.store)

			// Проверяем наличие или отсутствие ошибки
			if tc.wanterr {
				assert.Error(t, err, "Ожидалась ошибка, но её нет")
			} else {
				assert.NoError(t, err, "Ошибка не ожидалась, но она есть")
			}
		})
	}
}
