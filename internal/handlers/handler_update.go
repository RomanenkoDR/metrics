package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metric := chi.URLParam(r, "metric")
	value := chi.URLParam(r, "value")

	switch metricType {
	case counterType:
		v, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.Store.UpdateCounter(metric, storage.Counter(v))
	case gaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.Store.UpdateGauge(metric, storage.Gauge(v))
	default:
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
}

func (h *Handler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var m Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case counterType:
		if m.Delta == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.Store.UpdateCounter(m.ID, storage.Counter(*m.Delta))
		w.WriteHeader(http.StatusOK)
	case gaugeType:
		if m.Value == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.Store.UpdateGauge(m.ID, storage.Gauge(*m.Value))
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}
}

func (h *Handler) HandleUpdateBatch(w http.ResponseWriter, r *http.Request) {
	var m []Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, v := range m {
		switch v.MType {
		case counterType:
			if v.Delta == nil {
				http.Error(w, "metric value should not be empty", http.StatusBadRequest)
				return
			}
			h.Store.UpdateCounter(v.ID, storage.Counter(*v.Delta))
		case gaugeType:
			if v.Value == nil {
				http.Error(w, "metric value should not be empty", http.StatusBadRequest)
				return
			}
			h.Store.UpdateGauge(v.ID, storage.Gauge(*v.Value))
		default:
			http.Error(w, "Incorrect metric type", http.StatusBadRequest)
		}
	}
	w.WriteHeader(http.StatusOK)

}
