package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	MyTypeGauge   string = "gauge"
	MyTypeCounter string = "counter"
	PostUpdate    string = "/update"
	MType         string = "metricType"
	MName         string = "metricName"
	MValue        string = "metricValue"
)

// http://localhost:8080/update/counter/someMetric/527

// curl -v -X POST 'http://localhost:8080/update/counter/someMetric/527'

func gauge(res http.ResponseWriter, req *http.Request) {
	if MName == "" {
		fmt.Println("Metric name cannot be empty", http.StatusNotFound)
		return
	}
	if req.Method != http.MethodPost {
		res.Write([]byte("Метод отличен от POST"))
		fmt.Println(("Метод отличен от POST. Status code: "), http.StatusBadRequest)
		return
	}
	res.Write([]byte("Прошли в /update"))
	// fmt.Println(req.RequestURI)
	// fmt.Println(req.PathValue("metricType"))
	fmt.Println("Response OK. Status code: ", http.StatusOK)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Метод отличен от POST"))
		return
	}
	w.Write([]byte("Ответ"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(PostUpdate+MType+MName+MValue, gauge)
	mux.HandleFunc("/", mainPage)

	log.Println("Запуск веб-сервера на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
