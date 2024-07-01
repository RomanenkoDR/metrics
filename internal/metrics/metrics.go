package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"

	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

// Структура Metrics для хранения данных метрик
type Metrics struct {
	ID    string                `json:"id"`              // имя метрики
	MType string                `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta memStoragePcg.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value memStoragePcg.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// ReadMemStats считывает статистику памяти и обновляет метрики в хранилище
// Считывает статистику памяти с помощью runtime.ReadMemStats и
// обновляет значения метрик в хранилище.
// Также добавляет случайное значение и увеличивает счетчик запросов.
func ReadMemStats(m *memStoragePcg.MemStorage) {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	m.UpdateGauge("Alloc", memStoragePcg.Gauge(stat.Alloc))
	m.UpdateGauge("BuckHashSys", memStoragePcg.Gauge(stat.BuckHashSys))
	m.UpdateGauge("Frees", memStoragePcg.Gauge(stat.Frees))
	m.UpdateGauge("GCCPUFraction", memStoragePcg.Gauge(stat.GCCPUFraction))
	m.UpdateGauge("GCSys", memStoragePcg.Gauge(stat.GCSys))
	m.UpdateGauge("HeapAlloc", memStoragePcg.Gauge(stat.HeapAlloc))
	m.UpdateGauge("HeapIdle", memStoragePcg.Gauge(stat.HeapIdle))
	m.UpdateGauge("HeapInuse", memStoragePcg.Gauge(stat.HeapInuse))
	m.UpdateGauge("HeapObjects", memStoragePcg.Gauge(stat.HeapObjects))
	m.UpdateGauge("HeapReleased", memStoragePcg.Gauge(stat.HeapReleased))
	m.UpdateGauge("HeapSys", memStoragePcg.Gauge(stat.HeapSys))
	m.UpdateGauge("LastGC", memStoragePcg.Gauge(stat.LastGC))
	m.UpdateGauge("Lookups", memStoragePcg.Gauge(stat.Lookups))
	m.UpdateGauge("MCacheInuse", memStoragePcg.Gauge(stat.MCacheInuse))
	m.UpdateGauge("MCacheSys", memStoragePcg.Gauge(stat.MCacheSys))
	m.UpdateGauge("MSpanInuse", memStoragePcg.Gauge(stat.MSpanInuse))
	m.UpdateGauge("MSpanSys", memStoragePcg.Gauge(stat.MSpanSys))
	m.UpdateGauge("Mallocs", memStoragePcg.Gauge(stat.Mallocs))
	m.UpdateGauge("NextGC", memStoragePcg.Gauge(stat.NextGC))
	m.UpdateGauge("NumForcedGC", memStoragePcg.Gauge(stat.NumForcedGC))
	m.UpdateGauge("NumGC", memStoragePcg.Gauge(stat.NumGC))
	m.UpdateGauge("OtherSys", memStoragePcg.Gauge(stat.OtherSys))
	m.UpdateGauge("PauseTotalNs", memStoragePcg.Gauge(stat.PauseTotalNs))
	m.UpdateGauge("StackInuse", memStoragePcg.Gauge(stat.StackInuse))
	m.UpdateGauge("StackSys", memStoragePcg.Gauge(stat.StackSys))
	m.UpdateGauge("Sys", memStoragePcg.Gauge(stat.Sys))
	m.UpdateGauge("TotalAlloc", memStoragePcg.Gauge(stat.TotalAlloc))
	m.UpdateGauge("RandomValue", memStoragePcg.Gauge(rand.Float32()))
	m.UpdateCounter("PollCount", memStoragePcg.Counter(1))
}

// Compress сжимает входные данные data и возвращает сжатый результат или ошибку, если процесс сжатия не удался
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	// Создаем новый Writer с лучшей степенью сжатия
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		// Возвращаем ошибку, если не удалось инициализировать Writer для сжатия
		return nil, fmt.Errorf("не удалось инициализировать Writer для сжатия: %v", err)
	}

	// Пишем данные для сжатия во временный буфер
	_, err = w.Write(data)
	if err != nil {
		// Возвращаем ошибку, если не удалось записать данные во временный буфер для сжатия
		return nil, fmt.Errorf("не удалось записать данные во временный буфер для сжатия: %v", err)
	}

	// Закрываем Writer, завершая процесс сжатия
	err = w.Close()
	if err != nil {
		// Возвращаем ошибку, если не удалось завершить процесс сжатия
		return nil, fmt.Errorf("не удалось завершить процесс сжатия: %v", err)
	}

	// Возвращаем сжатые данные
	return b.Bytes(), nil
}

// PushMetricsToServer отправляет метрики на сервер
// Отправляет данные метрик на указанный сервер.
// Для каждой метрики формируется JSON-запрос, который затем отправляется на сервер.
// Если сервер возвращает ошибку, функция возвращает сообщение об ошибке.
func PushMetricsToServer(serverAddress string, m memStoragePcg.MemStorage) error {
	// Переменная для хранения метрик
	var metrics Metrics

	// Формирование URL для отправки запросов на сервер
	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	// отправка запроса на сервер для каждой метрики
	for k, v := range m.Data {
		switch v := v.(type) {
		case memStoragePcg.Gauge:
			metrics = Metrics{ID: k, MType: Gauge, Value: v}
		case memStoragePcg.Counter:
			metrics = Metrics{ID: k, MType: Counter, Delta: v}
		default:
			return fmt.Errorf("неизвестный тип метрики")
		}

		// Сериализация метрики в JSON
		data, err := json.Marshal(metrics)
		if err != nil {
			return err
		}

		// fmt.Println(string(data))
		// Сжатие данных
		data, err = Compress(data)
		if err != nil {
			return err
		}
		// Создание нового HTTP запроса
		request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		// Установка заголовков для сжатия и типа контента
		request.Header.Set("Content-Type", ContentTypeText)
		request.Header.Set("Content-Encoding", Compression)
		request.Header.Set("Accept-Encoding", Compression)

		// Выполнение HTTP запроса
		client := &http.Client{}
		resp, err := client.Do(request)

		if err != nil {
			return err
		}

		// Проверка статуса ответа
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("%s: %s; %s",
				"не удалось отправить метрики на сервер",
				resp.Status, b)
		}
		// Закрытие тела ответа
		defer resp.Body.Close()

	}
	return nil
}
