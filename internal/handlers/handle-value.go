package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handle URI request to return value
func (h *Handler) HandleValue(w http.ResponseWriter, r *http.Request) {
	metric := chi.URLParam(r, "metric")
	v, err := h.Store.Get(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Fprint(w, v)
}
