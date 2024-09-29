package handlers

import (
	"net/http"
	"strconv"

	memPcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// get context params
	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	// find out metric type
	switch metricType {
	case counterType:
		v, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.Store.UpdateCounter(metric, memPcg.Counter(v))
	case gaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.Store.UpdateGauge(metric, memPcg.Gauge(v))
	default:
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
}
