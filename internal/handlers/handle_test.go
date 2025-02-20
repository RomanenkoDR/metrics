package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// TestHandleMain проверяет основной обработчик ("/")
func TestHandleMain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		request    string
		httpMethod string
		wantCode   int
	}{
		{
			name:       "Root Page",
			request:    "/",
			httpMethod: http.MethodGet,
			wantCode:   http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.httpMethod, tc.request, nil)
			w := httptest.NewRecorder()

			h := NewHandler()
			h.HandleMain(w, req)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tc.wantCode, result.StatusCode)
			require.Equal(t, "text/html; charset=utf-8", result.Header.Get("Content-Type"))
		})
	}
}

// TestHandleUpdate проверяет обновление метрик
func TestHandleUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		wantCode    int
	}{
		{"Valid Gauge", "gauge", "m01", "1.3", http.StatusOK},
		{"Invalid Gauge", "gauge", "m01", "1ad3", http.StatusBadRequest},
		{"Valid Counter", "counter", "m02", "1", http.StatusOK},
		{"Invalid Counter", "counter", "m02", "1.4", http.StatusBadRequest},
		{"Invalid Metric Type", "nosuchmetric", "m02", "1.4", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqURL := "/update/" + tc.metricType + "/" + tc.metricName + "/" + tc.metricValue
			req := httptest.NewRequest(http.MethodPost, reqURL, nil)

			rCtx := chi.NewRouteContext()
			rCtx.URLParams.Add("type", tc.metricType)
			rCtx.URLParams.Add("metric", tc.metricName)
			rCtx.URLParams.Add("value", tc.metricValue)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

			h := NewHandler()
			h.HandleUpdate(w, req)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tc.wantCode, result.StatusCode)
		})
	}
}

// TestHandleValue проверяет получение значения метрики
func TestHandleValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		wantCode    int
	}{
		{"Existing Gauge", "gauge", "t1", "1.2", http.StatusOK},
		{"Existing Counter", "counter", "t2", "2", http.StatusOK},
		{"Nonexistent Metric", "counter", "t3", "3", http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler()

			// Записываем метрику перед тестом
			wUpdate := httptest.NewRecorder()
			updateURL := "/update/" + tc.metricType + "/" + tc.metricName + "/" + tc.metricValue
			reqUpdate := httptest.NewRequest(http.MethodPost, updateURL, nil)

			rCtx := chi.NewRouteContext()
			rCtx.URLParams.Add("type", tc.metricType)
			rCtx.URLParams.Add("metric", tc.metricName)
			rCtx.URLParams.Add("value", tc.metricValue)
			reqUpdate = reqUpdate.WithContext(context.WithValue(reqUpdate.Context(), chi.RouteCtxKey, rCtx))

			h.HandleUpdate(wUpdate, reqUpdate)

			// Тестируем получение метрики
			w := httptest.NewRecorder()
			valueURL := "/value/" + tc.metricType + "/" + tc.metricName
			req := httptest.NewRequest(http.MethodGet, valueURL, nil)

			rCtx = chi.NewRouteContext()
			rCtx.URLParams.Add("type", tc.metricType)
			rCtx.URLParams.Add("metric", tc.metricName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rCtx))

			h.HandleValue(w, req)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tc.wantCode, result.StatusCode)
		})
	}
}

// TestHandleUpdateBatch проверяет массовое обновление метрик через JSON
func TestHandleUpdateBatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		metrics  []Metrics
		wantCode int
	}{
		{
			name: "Valid Batch",
			metrics: []Metrics{
				{ID: "cpu_load", MType: "gauge", Value: ptrFloat64(0.75)},
				{ID: "requests", MType: "counter", Delta: ptrInt64(10)},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Invalid Batch - Wrong Type",
			metrics: []Metrics{
				{ID: "bad_metric", MType: "unknown", Value: ptrFloat64(1.23)},
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler()

			data, err := json.Marshal(tc.metrics)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/updates/", strings.NewReader(string(data)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.HandleUpdateBatch(w, req)

			result := w.Result()
			defer result.Body.Close()

			require.Equal(t, tc.wantCode, result.StatusCode)
		})
	}
}

// Вспомогательные функции для указателей
func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}
