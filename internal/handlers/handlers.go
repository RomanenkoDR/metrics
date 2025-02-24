package handlers

import (
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/jackc/pgx/v5"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Handler struct {
	Store          *storage.MemStorage
	DBconn         *pgx.Conn
	PrivateKeyPath string
}

const counterType = "counter"
const gaugeType = "gauge"

func NewHandler() Handler {
	return Handler{
		Store: storage.New(),
	}
}
