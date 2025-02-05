package main

import (
	"fmt"
	"github.com/RomanenkoDR/metrics/internal/config/server"
)
import _ "net/http/pprof"

var (
	buildVersion = "N/A" // VERSION 1-курс, 8-спринт, 20-инкремент
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	server.Run()
}
