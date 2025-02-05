package main

import (
	"github.com/RomanenkoDR/metrics/internal/analysers/noosexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
)

func main() {
	// Инициализируем анализаторы из пакета go/analysis/passes
	mychecks := []*analysis.Analyzer{
		inspect.Analyzer,
		shadow.Analyzer,
		printf.Analyzer,
	}

	// Добавляем анализаторы из пакета staticcheck.io
	staticcheckAnalyzers := staticcheck.Analyzers
	for _, v := range staticcheckAnalyzers {
		if v.Analyzer != nil {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// Добавляем собственный анализатор
	mychecks = append(mychecks, noosexit.Analyzer)

	// Запускаем multichecker
	multichecker.Main(mychecks...)
}
