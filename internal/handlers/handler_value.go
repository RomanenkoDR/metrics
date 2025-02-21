package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// HandleValue URI request to return value
func (h *Handler) HandleValue(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	v, err := h.Store.Get(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Fprint(w, v)
}

// HandleValueJSON request to return value
func (h *Handler) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	var m Metrics

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case counterType:
		v, ok := h.Store.CounterData[m.ID]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		vPtr := int64(v)
		m.Delta = &vPtr
	case gaugeType:
		v, ok := h.Store.GaugeData[m.ID]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		vPtr := float64(v)
		m.Value = &vPtr
	}

	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// respond to agent
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
