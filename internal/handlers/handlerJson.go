package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

func (h *handler) PostValueByJSON(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(m)

	if _, ok := h.store.Data[m.ID]; !ok {
		http.Error(w, "not found", http.StatusNotFound)

		return
	}

	switch m.MType {
	case Counter:
		v := int64(h.store.Data[m.ID].(memStoragePcg.Counter))
		m.Delta = &v
	case Gauge:
		v := float64(h.store.Data[m.ID].(memStoragePcg.Gauge))
		m.Value = &v
	}

	resp, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) PostUpdateJSON(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(m)

	switch m.MType {
	case Counter:
		if m.Delta == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.store.UpdateCounter(m.ID, memStoragePcg.Counter(*m.Delta))
		w.WriteHeader(http.StatusOK)
	case Gauge:
		if m.Value == nil {
			http.Error(w, "metric value should not be empty", http.StatusBadRequest)
			return
		}
		h.store.UpdateGauge(m.ID, memStoragePcg.Gauge(*m.Value))
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Incorrect metric type", http.StatusBadRequest)
	}

}
