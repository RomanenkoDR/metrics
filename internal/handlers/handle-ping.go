package handlers

import (
	"context"
	"fmt"
	"net/http"
)

func (h *Handler) HandlePing(w http.ResponseWriter, r *http.Request) {
	v := "pong\n"
	err := dbconn.Ping(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, v)
}
