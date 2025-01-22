package main

import "github.com/RomanenkoDR/metrics/internal/config/server"
import _ "net/http/pprof"

func main() {
	server.Run()
}
