package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Функция для генерации подписи для теста
func generateSignature(key string, body []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// Тестирование CheckReqSign с правильной подписью
func TestCheckReqSign_ValidSignature(t *testing.T) {
	key := "test-secret-key"
	body := []byte("test data")

	// Генерация правильной подписи
	expectedSignature := generateSignature(key, body)

	// Создаем запрос с правильной подписью в заголовке
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("HashSHA256", expectedSignature)

	// Создаем рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик, который просто возвращает статус 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware CheckReqSign
	CheckReqSign(key)(handler).ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что заголовок HashSHA256 был добавлен в ответ
	assert.NotEmpty(t, rr.Header().Get("HashSHA256"))
}

// Тестирование CheckReqSign с неправильной подписью
func TestCheckReqSign_InvalidSignature(t *testing.T) {
	key := "test-secret-key"
	body := []byte("test data")

	// Генерация неправильной подписи (изменим тело запроса)
	incorrectBody := []byte("modified data")
	incorrectSignature := generateSignature(key, incorrectBody)

	// Создаем запрос с неправильной подписью в заголовке
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("HashSHA256", incorrectSignature)

	// Создаем рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик, который просто возвращает статус 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware CheckReqSign
	CheckReqSign(key)(handler).ServeHTTP(rr, req)

	// Проверяем, что статус код 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Проверяем, что в теле ответа присутствует ошибка о поврежденной подписи
	assert.Contains(t, rr.Body.String(), "Corrupted sign on request")
}

// Тестирование CheckReqSign без подписи
func TestCheckReqSign_NoSignature(t *testing.T) {
	key := "test-secret-key"
	body := []byte("test data")

	// Создаем запрос без подписи в заголовке
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	// Создаем рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик, который просто возвращает статус 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применяем middleware CheckReqSign
	CheckReqSign(key)(handler).ServeHTTP(rr, req)

	// Проверяем, что статус код 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что заголовок HashSHA256 был добавлен в ответ
	assert.NotEmpty(t, rr.Header().Get("HashSHA256"))
}
