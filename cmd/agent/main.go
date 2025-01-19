package main

import (
	"github.com/RomanenkoDR/metrics/internal/config/agent"
	_ "net/http/pprof"
)

func main() {
	agent.RunAgent()
}
