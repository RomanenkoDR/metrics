package main

import (
	"log"
	"time"

	agentConfigPcg "github.com/RomanenkoDR/metrics/internal/config/agentConfig"
	metricsPcg "github.com/RomanenkoDR/metrics/internal/metrics"
	memStoragePcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

func main() {
	//parse cli options
	config, err := agentConfigPcg.ParseOptions()
	if err != nil {
		panic(err)
	}

	// initiate tickers
	pollTicker := time.NewTicker(time.Second * time.Duration(config.PollInterval))
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Second * time.Duration(config.ReportInterval))
	defer reportTicker.Stop()

	//initiate new storage
	m := memStoragePcg.New()

	//collect data from MemStats and send to the server
	for {
		select {
		case <-pollTicker.C:
			metricsPcg.ReadMemStats(&m)
		case <-reportTicker.C:
			err := metricsPcg.ProcessReport(config.ServerAddress, m)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
