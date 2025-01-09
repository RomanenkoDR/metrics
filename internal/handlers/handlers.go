package handlers

import (
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

const counterType = "counter"
const gaugeType = "gauge"

// Create new handler and previous reports info from file it needed
func NewHandler() Handler {
	var h Handler
	h.Store = storage.New()

	return h
}

type Handler struct {
	Store  storage.MemStorage
	DBconn *pgx.Conn
}
