package main

import (
	"github.com/RomanenkoDR/metrics/internal/config/server"
	_ "net/http/pprof"
)

func main() {
	server.RunServer()
}
