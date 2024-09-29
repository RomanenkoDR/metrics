package main

import (
	"log"
	"time"

	cnfAgentPcg "github.com/RomanenkoDR/metrics/internal/config/agentcfg"
	metricPcg "github.com/RomanenkoDR/metrics/internal/metrics"
	memPcg "github.com/RomanenkoDR/metrics/internal/storage/mem"
)

func main() {
	//parse cli options
	cfg, err := cnfAgentPcg.ParseOptions()
	if err != nil {
		panic(err)
	}

	// initiate tickers
	pollTicker := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	defer reportTicker.Stop()

	//initiate new storage
	m := memPcg.New()

	//collect data from MemStats and send to the server
	for {
		select {
		case <-pollTicker.C:
			metricPcg.ReadMemStats(&m)
		case <-reportTicker.C:
			err := metricPcg.ProcessReport(cfg.ServerAddress, m)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
