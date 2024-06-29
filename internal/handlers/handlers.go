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

func (h *handler) GetListAllMetrics(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("../../internal/template/listMetricsPage.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := h.store.GetAll()
	metricsData := MetricsData{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	for k, v := range list {
		switch value := v.(type) {
		case memStoragePcg.Counter:
			metricsData.Counters[k] = int64(value)
		case memStoragePcg.Gauge:
			metricsData.Gauges[k] = float64(value)
		default:
			fmt.Printf("Unexpected type for key %s: %T\n", k, v)
		}
	}

	if err := tmpl.Execute(w, metricsData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) PostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	// get context params
	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	// find out metric type
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
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
}

func (h *handler) GetValueByName(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	v, err := h.store.Get(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Fprint(w, v)
}
