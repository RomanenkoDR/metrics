package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/RomanenkoDR/metrics/internal/storage"
)

type MetricsData struct {
	Gauges   map[string]storage.Gauge
	Counters map[string]storage.Counter
}

// HandleMain рендерит HTML-шаблон с метриками.
func (h *Handler) HandleMain(w http.ResponseWriter, r *http.Request) {
	// Путь к HTML-шаблону
	log.Println("запуска main функции")
	templatePath := "template/template.html"

	// Загрузка шаблона
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	// Подготовка данных для рендера
	data := MetricsData{
		Gauges:   h.Store.GetAllGauge(),
		Counters: h.Store.GetAllCounters(),
	}

	// Установка заголовков и рендер шаблона
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		return
	}
}
