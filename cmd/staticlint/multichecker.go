package main

//func main() {
//	// Инициализируем анализаторы из пакета go/analysis/passes
//	mychecks := []*analysis.Analyzer{
//		inspect.Analyzer,
//		shadow.Analyzer,
//		printf.Analyzer,
//	}
//
//	// Добавляем анализаторы из пакета staticcheck.io
//	staticcheckAnalyzers := staticcheck.Analyzers
//	for _, v := range staticcheckAnalyzers {
//		if v.Analyzer != nil {
//			mychecks = append(mychecks, v.Analyzer)
//		}
//	}
//
//	// Добавляем собственный анализатор
//	mychecks = append(mychecks, noosexit.Analyzer)
//
//	// Запускаем multichecker
//	multichecker.Main(mychecks...)
//}
