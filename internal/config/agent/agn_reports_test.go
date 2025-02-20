package agent

//import (
//	"context"
//	"encoding/json"
//	"github.com/RomanenkoDR/metrics/internal/storage"
//	"github.com/stretchr/testify/require"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)

//// Моковый обработчик сервера
//func mockServer(t *testing.T, expectedStatus int, expectedResponse string) *httptest.Server {
//	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if expectedStatus != http.StatusOK {
//			http.Error(w, expectedResponse, expectedStatus)
//			return
//		}
//
//		// Читаем тело запроса
//		body, err := io.ReadAll(r.Body)
//		require.NoError(t, err)
//
//		// Логируем полученные данные
//		t.Logf("Полученные данные: %s", string(body))
//
//		w.WriteHeader(http.StatusOK)
//	}))
//}
//
//func TestProcessReport(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name        string
//		cryptoKey   string
//		storage     storage.MemStorage
//		wantErr     bool
//		wantErrText string
//	}{
//		{
//			name:      "Valid request without encryption",
//			cryptoKey: "",
//			storage: storage.MemStorage{
//				CounterData: map[string]storage.Counter{
//					"requests": 10,
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:      "Valid request with encryption",
//			cryptoKey: "/path/to/public.pem",
//			storage: storage.MemStorage{
//				GaugeData: map[string]storage.Gauge{
//					"cpu_load": 0.75,
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:        "Invalid metric type",
//			cryptoKey:   "",
//			storage:     storage.MemStorage{},
//			wantErr:     true,
//			wantErrText: "can't send report to the server: 400 Bad Request",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			server := mockServer(t, http.StatusOK, "OK")
//			defer server.Close()
//
//			err := ProcessReport(strings.TrimPrefix(server.URL, "http://"), tc.cryptoKey, tc.storage)
//
//			if tc.wantErr {
//				require.Error(t, err)
//				require.Contains(t, err.Error(), tc.wantErrText)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestProcessBatch(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name        string
//		cryptoKey   string
//		storage     storage.MemStorage
//		wantErr     bool
//		wantErrText string
//	}{
//		{
//			name:      "Valid batch without encryption",
//			cryptoKey: "",
//			storage: storage.MemStorage{
//				CounterData: map[string]storage.Counter{
//					"requests": 20,
//				},
//				GaugeData: map[string]storage.Gauge{
//					"cpu_load": 0.95,
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:      "Valid batch with encryption",
//			cryptoKey: "/path/to/public.pem",
//			storage: storage.MemStorage{
//				CounterData: map[string]storage.Counter{
//					"errors": 5,
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name:        "Empty batch request",
//			cryptoKey:   "",
//			storage:     storage.MemStorage{},
//			wantErr:     true,
//			wantErrText: "can't send report to the server: 400 Bad Request",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			server := mockServer(t, http.StatusOK, "OK")
//			defer server.Close()
//
//			err := ProcessBatch(context.Background(), strings.TrimPrefix(server.URL, "http://"), tc.cryptoKey, tc.storage)
//
//			if tc.wantErr {
//				require.Error(t, err)
//				require.Contains(t, err.Error(), tc.wantErrText)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestSendRequest_ErrorCases(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name          string
//		serverStatus  int
//		serverResp    string
//		cryptoKey     string
//		inputData     interface{}
//		expectErr     bool
//		expectedError string
//	}{
//		{
//			name:          "Server returns 400",
//			serverStatus:  http.StatusBadRequest,
//			serverResp:    "Bad Request",
//			cryptoKey:     "",
//			inputData:     map[string]string{"id": "test", "value": "100"},
//			expectErr:     true,
//			expectedError: "can't send report to the server: 400 Bad Request",
//		},
//		{
//			name:          "Invalid JSON Serialization",
//			serverStatus:  http.StatusOK,
//			serverResp:    "OK",
//			cryptoKey:     "",
//			inputData:     make(chan int), // Канал не сериализуем в JSON
//			expectErr:     true,
//			expectedError: "json: unsupported type: chan int",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			server := mockServer(t, tc.serverStatus, tc.serverResp)
//			defer server.Close()
//
//			data, err := json.Marshal(tc.inputData)
//			if err != nil && tc.expectErr {
//				require.Error(t, err)
//				require.Contains(t, err.Error(), tc.expectedError)
//				return
//			}
//
//			err = sendRequest(strings.TrimPrefix(server.URL, "http://"), data, tc.cryptoKey)
//
//			if tc.expectErr {
//				require.Error(t, err)
//				require.Contains(t, err.Error(), tc.expectedError)
//			} else {
//				require.NoError(t, err)
//			}
//		})
//	}
//}
