package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестирование GzipHandle для сжатия контента
func TestGzipHandleCompress(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("Hello, Gzip"))
	})

	// Создаем запрос с заголовком Accept-Encoding, чтобы указать, что клиент поддерживает gzip
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	// Создаем рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Выполняем обработку запроса с gzip-обработчиком
	GzipHandle(handler).ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что в заголовках установлен Content-Encoding: gzip
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Проверяем, что тело ответа сжато
	body := rr.Body.Bytes()
	if !isGzipped(body) {
		t.Errorf("Expected gzipped content, but got: %v", body)
	}

	// Распаковываем и проверяем содержимое
	gzReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	uncompressedBody, err := io.ReadAll(gzReader)
	if err != nil {
		t.Fatalf("Failed to read gzipped body: %v", err)
	}

	// Проверяем содержимое после распаковки
	assert.Equal(t, "Hello, Gzip", string(uncompressedBody))
}

// Тестирование GzipHandle без сжатия
func TestGzipHandleNoCompress(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("Hello, World"))
	})

	// Создаем запрос без заголовка Accept-Encoding
	req := httptest.NewRequest("GET", "/", nil)

	// Создаем рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Выполняем обработку запроса с gzip-обработчиком
	GzipHandle(handler).ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что в заголовках нет Content-Encoding
	assert.Empty(t, rr.Header().Get("Content-Encoding"))

	// Проверяем тело ответа
	assert.Equal(t, "Hello, World", rr.Body.String())
}

// Проверка, является ли тело gzipped
func isGzipped(data []byte) bool {
	return len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b
}
