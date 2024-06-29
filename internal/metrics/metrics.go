package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"

	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

type Metrics struct {
	ID    string                `json:"id"`    // имя метрики
	MType string                `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta memStoragePcg.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value memStoragePcg.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}

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

func ProcessReport(serverAddress string, m memStoragePcg.MemStorage) error {
	// metric type variable

	var metrics Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, v := range m.Data {
		switch v := v.(type) {
		case memStoragePcg.Gauge:
			metrics = Metrics{ID: k, MType: Gauge, Value: v}
		case memStoragePcg.Counter:
			metrics = Metrics{ID: k, MType: Counter, Delta: v}
		default:
			return fmt.Errorf("uknown type of metric")
		}

		data, err := json.Marshal(metrics)
		if err != nil {
			return err
		}

		//         fmt.Println(string(data))

		request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", ContentType)

		client := &http.Client{}
		resp, err := client.Do(request)

		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("%s: %s; %s",
				"Can't send report to the server",
				resp.Status,
				b)
		}

		defer resp.Body.Close()

	}
	return nil
}
