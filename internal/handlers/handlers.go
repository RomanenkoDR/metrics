package handlers

import (
	"encoding/json"
	"log"
	"os"

	memPcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
	"github.com/jackc/pgx/v5"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Handler struct {
	Store memPcg.MemStorage
}

const counterType = "counter"
const gaugeType = "gauge"

var dbconn *pgx.Conn

func SetDB(db *pgx.Conn) {
	dbconn = db
}

// Create new handler and previous reports info from file it needed
func NewHandler(filename string, restore bool) (Handler, error) {
	var h Handler
	h.Store = memPcg.New()

	// Read saved metrics from file
	if restore {
		f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return h, err
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		err = decoder.Decode(&h.Store)
		if err != nil {
			log.Println("Could not restore data", err)
		}
	}
	return h, nil
}
