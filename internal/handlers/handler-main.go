package handlers

import (
	"fmt"
	"net/http"
)

// type Data struct {
// 	Gauges   map[string]float64
// 	Counters map[string]int64
// }

// func HandleMainTemplate(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	_, err := res.Write([]byte(""))
// 	if err != nil {
// 		log.Printf("Failed to write response: %v", err)
// 		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
// 	}

// 	if req.URL.Path != "/" {
// 		res.WriteHeader(http.StatusNotFound)
// 		res.Write([]byte("Такого нет( - 404)"))
// 		return
// 	}

// 	if req.Method == http.MethodGet {
// 		gauges :
// 		counters := memPcg.Data.GetAllCounters()

// 		data := Data{
// 			Gauges:   gauges,
// 			Counters: counters,
// 		}

// 		tmplPath := "../../internal/templates/metrics.html"
// 		tmpl, err := template.ParseFiles(tmplPath)
// 		if err != nil {
// 			http.Error(res, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		res.WriteHeader(http.StatusOK)

// 		err = tmpl.Execute(res, data)
// 		if err != nil {
// 			http.Error(res, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
// 		}
// 	}
// }

func (h *Handler) HandleMain(w http.ResponseWriter, r *http.Request) {
	//write static html page with all the items to the response; unsorted
	body := `
        <!DOCTYPE html>
        <html>
            <head>
                <title>All tuples</title>
            </head>
            <body>
            <table>
                <tr>
                    <td>Metric</td>
                    <td>Value</td>
                </tr>
    `
	listC := h.Store.GetAllCounters()
	for k, v := range listC {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	listG := h.Store.GetAllGauges()
	for k, v := range listG {
		body = body + fmt.Sprintf("<tr>\n<td>%s</td>\n", k)
		body = body + fmt.Sprintf("<td>%v</td>\n</tr>\n", v)
	}

	body = body + " </table>\n </body>\n</html>"

	// respond to agent
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
