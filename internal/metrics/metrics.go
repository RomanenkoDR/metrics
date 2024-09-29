package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strings"

	memPcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

type Metrics struct {
	ID    string         `json:"id"`    // имя метрики
	MType string         `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta memPcg.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value memPcg.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}

const contentType string = "application/json"
const compression string = "gzip"

const counterType string = "counter"
const gaugeType string = "gauge"

func ReadMemStats(m *memPcg.MemStorage) {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	m.UpdateGauge("Alloc", memPcg.Gauge(stat.Alloc))
	m.UpdateGauge("BuckHashSys", memPcg.Gauge(stat.BuckHashSys))
	m.UpdateGauge("Frees", memPcg.Gauge(stat.Frees))
	m.UpdateGauge("GCCPUFraction", memPcg.Gauge(stat.GCCPUFraction))
	m.UpdateGauge("GCSys", memPcg.Gauge(stat.GCSys))
	m.UpdateGauge("HeapAlloc", memPcg.Gauge(stat.HeapAlloc))
	m.UpdateGauge("HeapIdle", memPcg.Gauge(stat.HeapIdle))
	m.UpdateGauge("HeapInuse", memPcg.Gauge(stat.HeapInuse))
	m.UpdateGauge("HeapObjects", memPcg.Gauge(stat.HeapObjects))
	m.UpdateGauge("HeapReleased", memPcg.Gauge(stat.HeapReleased))
	m.UpdateGauge("HeapSys", memPcg.Gauge(stat.HeapSys))
	m.UpdateGauge("LastGC", memPcg.Gauge(stat.LastGC))
	m.UpdateGauge("Lookups", memPcg.Gauge(stat.Lookups))
	m.UpdateGauge("MCacheInuse", memPcg.Gauge(stat.MCacheInuse))
	m.UpdateGauge("MCacheSys", memPcg.Gauge(stat.MCacheSys))
	m.UpdateGauge("MSpanInuse", memPcg.Gauge(stat.MSpanInuse))
	m.UpdateGauge("MSpanSys", memPcg.Gauge(stat.MSpanSys))
	m.UpdateGauge("Mallocs", memPcg.Gauge(stat.Mallocs))
	m.UpdateGauge("NextGC", memPcg.Gauge(stat.NextGC))
	m.UpdateGauge("NumForcedGC", memPcg.Gauge(stat.NumForcedGC))
	m.UpdateGauge("NumGC", memPcg.Gauge(stat.NumGC))
	m.UpdateGauge("OtherSys", memPcg.Gauge(stat.OtherSys))
	m.UpdateGauge("PauseTotalNs", memPcg.Gauge(stat.PauseTotalNs))
	m.UpdateGauge("StackInuse", memPcg.Gauge(stat.StackInuse))
	m.UpdateGauge("StackSys", memPcg.Gauge(stat.StackSys))
	m.UpdateGauge("Sys", memPcg.Gauge(stat.Sys))
	m.UpdateGauge("TotalAlloc", memPcg.Gauge(stat.TotalAlloc))
	m.UpdateGauge("RandomValue", memPcg.Gauge(rand.Float64()))
	m.UpdateCounter("PollCount", memPcg.Counter(1))
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func sendReport(serverAddress string, metrics Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	data, err = compress(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Content-Encoding", compression)
	request.Header.Set("Accept-Encoding", compression)

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
	return nil
}

func ProcessReport(serverAddress string, m memPcg.MemStorage) error {

	var metrics Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, v := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: v}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	for k, v := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: v}
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	return nil
}
