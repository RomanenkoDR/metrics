package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	s "github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestProcessBatch(t *testing.T) {
	t.Parallel()

	// Моковый ответ сервера
	responseBody := "response"

	tests := []struct {
		name         string
		store        s.MemStorage
		cryptoKey    string // путь к ключу шифрования
		wantErr      string // ожидаемая строка ошибки
		wantHTTPCode int
	}{
		{
			name: "Valid Batch request with gauge metric",
			store: s.MemStorage{
				GaugeData: map[string]s.Gauge{
					"metric1": 2.32,
					"metric2": 3.45,
				},
			},
			cryptoKey:    "", // Без шифрования
			wantErr:      "",
			wantHTTPCode: http.StatusOK,
		},
		{
			name: "Valid Batch request with encryption",
			store: s.MemStorage{
				GaugeData: map[string]s.Gauge{
					"metric1": 10.5,
				},
			},
			cryptoKey:    "/path/to/public.pem", // Шифрование включено
			wantErr:      "",
			wantHTTPCode: http.StatusOK,
		},
		{
			name:         "Empty Batch request",
			store:        s.MemStorage{CounterData: map[string]s.Counter{}},
			cryptoKey:    "",
			wantErr:      "",
			wantHTTPCode: http.StatusBadRequest,
		},
		{
			name: "Invalid Batch request with counter metric",
			store: s.MemStorage{
				CounterData: map[string]s.Counter{
					"counter1": 5,
				},
			},
			cryptoKey:    "",
			wantErr:      fmt.Sprintf("Can't send report to the server: 400 Bad Request; %s", responseBody),
			wantHTTPCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Создаём мок-сервер
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				require.Equal(t, "application/json", req.Header.Get("Content-Type")) // Проверяем заголовки
				http.Error(rw, responseBody, tc.wantHTTPCode)
			}))
			defer server.Close()

			// Отправка батча на мок-сервер
			ctx := context.Background()
			err := ProcessBatch(ctx, strings.TrimPrefix(server.URL, "http://"), tc.cryptoKey, tc.store)

			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, strings.TrimSpace(tc.wantErr))
			}
		})
	}
}
