package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

// PostValueByJSON обрабатывает запросы на получение значения метрики в формате JSON
// Читает запрос JSON, десериализует его и получает метрику из хранилища.
// В зависимости от типа метрики (Counter или Gauge), добавляет значение к ответу JSON и отправляет его клиенту.
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
		http.Error(w, "не найдено", http.StatusNotFound)

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

// PostUpdateJSON обрабатывает запросы на обновление метрик в формате JSON
// Читает запрос JSON, десериализует его и обновляет значение метрики в хранилище
// в зависимости от типа метрики (Counter или Gauge).
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
			http.Error(w, "Значение метрики не должно быть пустым", http.StatusBadRequest)
			return
		}
		h.store.UpdateCounter(m.ID, memStoragePcg.Counter(*m.Delta))
		w.WriteHeader(http.StatusOK)
	case Gauge:
		if m.Value == nil {
			http.Error(w, "Значение метрики не должно быть пустым", http.StatusBadRequest)
			return
		}
		h.store.UpdateGauge(m.ID, memStoragePcg.Gauge(*m.Value))
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Неверный тип метрики", http.StatusBadRequest)
	}

}
