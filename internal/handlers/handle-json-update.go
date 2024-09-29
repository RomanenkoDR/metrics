package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	storage "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

// Handle JSON request to update metric value
func (h *Handler) HandleJSONUpdate(w http.ResponseWriter, r *http.Request) {
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
