package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

// Функция для сжатия данных с использованием gzip
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	// Пишем исходные данные в gzip writer для сжатия
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	// Закрываем writer и завершаем процесс сжатия
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	// Возвращаем сжатые данные в виде байтового среза
	return b.Bytes(), nil
}
