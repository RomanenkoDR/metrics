package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Metrics struct {
	ID    string          `json:"id"`    // имя метрики
	MType string          `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta storage.Counter `json:"delta"` // значение метрики в случае передачи counter
	Value storage.Gauge   `json:"value"` // значение метрики в случае передачи gauge
}

const contentType string = "application/json"
const compression string = "gzip"

const counterType string = "counter"
const gaugeType string = "gauge"

// Renew metrics through runtime package
func ReadMemStats(m *storage.MemStorage, metricsCh chan storage.MemStorage) {
	var stat runtime.MemStats
	var mu sync.RWMutex

	mu.Lock()

	runtime.ReadMemStats(&stat)
	m.UpdateGauge("Alloc", storage.Gauge(stat.Alloc))
	m.UpdateGauge("BuckHashSys", storage.Gauge(stat.BuckHashSys))
	m.UpdateGauge("Frees", storage.Gauge(stat.Frees))
	m.UpdateGauge("GCCPUFraction", storage.Gauge(stat.GCCPUFraction))
	m.UpdateGauge("GCSys", storage.Gauge(stat.GCSys))
	m.UpdateGauge("HeapAlloc", storage.Gauge(stat.HeapAlloc))
	m.UpdateGauge("HeapIdle", storage.Gauge(stat.HeapIdle))
	m.UpdateGauge("HeapInuse", storage.Gauge(stat.HeapInuse))
	m.UpdateGauge("HeapObjects", storage.Gauge(stat.HeapObjects))
	m.UpdateGauge("HeapReleased", storage.Gauge(stat.HeapReleased))
	m.UpdateGauge("HeapSys", storage.Gauge(stat.HeapSys))
	m.UpdateGauge("LastGC", storage.Gauge(stat.LastGC))
	m.UpdateGauge("Lookups", storage.Gauge(stat.Lookups))
	m.UpdateGauge("MCacheInuse", storage.Gauge(stat.MCacheInuse))
	m.UpdateGauge("MCacheSys", storage.Gauge(stat.MCacheSys))
	m.UpdateGauge("MSpanInuse", storage.Gauge(stat.MSpanInuse))
	m.UpdateGauge("MSpanSys", storage.Gauge(stat.MSpanSys))
	m.UpdateGauge("Mallocs", storage.Gauge(stat.Mallocs))
	m.UpdateGauge("NextGC", storage.Gauge(stat.NextGC))
	m.UpdateGauge("NumForcedGC", storage.Gauge(stat.NumForcedGC))
	m.UpdateGauge("NumGC", storage.Gauge(stat.NumGC))
	m.UpdateGauge("OtherSys", storage.Gauge(stat.OtherSys))
	m.UpdateGauge("PauseTotalNs", storage.Gauge(stat.PauseTotalNs))
	m.UpdateGauge("StackInuse", storage.Gauge(stat.StackInuse))
	m.UpdateGauge("StackSys", storage.Gauge(stat.StackSys))
	m.UpdateGauge("Sys", storage.Gauge(stat.Sys))
	m.UpdateGauge("TotalAlloc", storage.Gauge(stat.TotalAlloc))
	m.UpdateGauge("RandomValue", storage.Gauge(rand.Float32()))
	m.UpdateCounter("PollCount", storage.Counter(1))

	// gopsutil metrics
	vmem, _ := mem.VirtualMemory()
	cpu1, _ := cpu.Percent(time.Duration(0), true)

	m.UpdateGauge("TotalMemory", storage.Gauge(vmem.Total))
	m.UpdateGauge("FreeMemory", storage.Gauge(vmem.Free))
	m.UpdateGauge("CPUutilization1", storage.Gauge(cpu1[0]))

	mu.Unlock()

	metricsCh <- *m
}

// Compress function profides fast compression
// for requests to send to the server
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %vmem", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %vmem", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %vmem", err)
	}
	return b.Bytes(), nil
}

// Send atomic metric report to the server
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

// Process all the metrics and send them to the server one by one
func ProcessReport(serverAddress string, m storage.MemStorage) error {
	// metric type variable

	var metrics Metrics

	serverAddress = strings.Join([]string{"http:/", serverAddress, "update/"}, "/")

	//send request to the server
	for k, vmem := range m.CounterData {
		metrics = Metrics{ID: k, MType: counterType, Delta: vmem}
		log.Println(metrics)
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}

	for k, vmem := range m.GaugeData {
		metrics = Metrics{ID: k, MType: gaugeType, Value: vmem}
		err := sendReport(serverAddress, metrics)
		if err != nil {
			return err
		}
	}
	return nil
}
