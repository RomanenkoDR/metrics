package db

import (
	"github.com/RomanenkoDR/metrics/internal/db/db_types"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"time"
)

// Save сохраняет данные в базе данных с указанным интервалом.
func Save(db *db_types.Database, interval int, s storage.MemStorage) error {
	time.Sleep(time.Second * time.Duration(interval))
	return Write(db, s)
}
