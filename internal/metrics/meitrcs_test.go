package metrics

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	memPcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
// 	"github.com/go-playground/assert"
// )

// func TestProcessReport(t *testing.T) {
// 	// http server response body
// 	responseBody := "response"

// 	tests := []struct {
// 		name     string
// 		store    memPcg.MemStorage
// 		wanterr  error
// 		wantcode int
// 	}{
// 		{
// 			name: "Test Valid Post request gauge metric",
// 			store: memPcg.MemStorage{
// 				GaugeData: map[string]memPcg.Gauge{
// 					"valid": memPcg.Gauge(2.32),
// 				},
// 			},
// 			wanterr:  nil,
// 			wantcode: http.StatusOK,
// 		},
// 		{
// 			name:  "Test Empty metric",
// 			store: memPcg.MemStorage{CounterData: map[string]memPcg.Counter{}},
// 			// adding new line into format string as http server do
// 			wanterr:  nil,
// 			wantcode: http.StatusBadRequest,
// 		},
// 		{
// 			name: "Test Invalid Post request counter metric",
// 			store: memPcg.MemStorage{
// 				CounterData: map[string]memPcg.Counter{
// 					"valid": memPcg.Counter(2),
// 				},
// 			},
// 			// adding new line into format string as http server do
// 			wanterr: fmt.Errorf("%s: %s; %s\n",
// 				"Can't send report to the server",
// 				"400 Bad Request",
// 				responseBody),
// 			wantcode: http.StatusBadRequest,
// 		},
// 	}
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 				http.Error(rw, responseBody, tc.wantcode)
// 			}))
// 			defer server.Close()

// 			//             .serverAddress =
// 			err := ProcessReport(strings.Replace(server.URL, "http://", "", 1), tc.store)
// 			assert.Equal(t, tc.wanterr, err)
// 		})
// 	}
// }
