package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	store memStoragePcg.MemStorage
}

func NewHandler() handler {
	return handler{store: memStoragePcg.New()}
}

type MetricsData struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

// GetListAllMetrics обрабатывает запросы на получение списка всех метрик
// Загружает HTML-шаблон для отображения всех метрик.
// Получает все метрики из хранилища и сортирует их по типам (Counter и Gauge).
// Выполняет шаблон с данными метрик и отправляет результат клиенту.
func (h *handler) GetListAllMetrics(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("../../internal/template/listMetricsPage.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получение всех метрик из хранилища
	list := h.store.GetAll()
	metricsData := MetricsData{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	// Заполнение данных метрик
	for k, v := range list {
		switch value := v.(type) {
		case memStoragePcg.Counter:
			metricsData.Counters[k] = int64(value)
		case memStoragePcg.Gauge:
			metricsData.Gauges[k] = float64(value)
		default:
			fmt.Printf("Неожиданный тип для ключа %s: %T\n", k, v)
		}
	}

	if err := tmpl.Execute(w, metricsData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// // PostUpdateMetric обрабатывает запросы на обновление метрик
// Получает параметры метрики (тип, имя и значение) из URL.
// В зависимости от типа метрики (Counter или Gauge), обновляет значение метрики в хранилище.
func (h *handler) PostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	// Определение типа метрики
	switch metricType {
	case Counter:
		v, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.store.UpdateCounter(metric, memStoragePcg.Counter(v))
	case Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.store.UpdateGauge(metric, memStoragePcg.Gauge(v))
	default:
		http.Error(w, "Неверный тип метрики", http.StatusBadRequest)
	}
}

// GetValueByName обрабатывает запросы на получение значения метрики по имени
// Получает значение метрики по имени из URL и отправляет его клиенту.
func (h *handler) GetValueByName(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	v, err := h.store.Get(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Fprint(w, v)
}
