package handlers

import (
	"fmt"
	"net/http"
)

// Return list with all the metrics

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

	listG := h.Store.GetAllGauge()
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
